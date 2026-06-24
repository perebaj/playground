# Batch Size 400 Error Reproduction

Reproduces 400 errors caused by a partner's OTel Collector config that batches 100K log records into a single HTTP export request.

## Hypothesis

The `batch/logs_export` processor with `send_batch_size: 100000` creates OTLP HTTP requests too large for raincatcher to accept, resulting in HTTP 400 (InvalidArgument).

The serialized protobuf payload exceeds the NATS limits in the pipeline:

| Limit | Dev | Prod |
|-------|-----|------|
| NATS `max_payload` | 1 MB | 8 MB |
| `og-logs` stream `max_msg_size` | 1 MiB | 2 MiB |
| gRPC max recv (raincatcher) | 4 MB | 4 MB |

## Setup

```bash
cd repro-batch-400
export OLLY_GARDEN_API_KEY="og_sk_..."
```

## Test 1: Large batch (100K) — reproduce the 400 errors

```bash
docker compose -f docker-compose-large-batch.yml up
```

Watch for errors:

```bash
docker logs -f otel-collector-large-batch 2>&1 | grep -iE "400|error|drop"
```

Stop:

```bash
docker compose -f docker-compose-large-batch.yml down
```

## Test 2: Small batch (5K) — should succeed

```bash
docker compose -f docker-compose-small-batch.yml up
```

Watch for errors:

```bash
docker logs -f otel-collector-small-batch 2>&1 | grep -iE "400|error|drop"
```

Stop:

```bash
docker compose -f docker-compose-small-batch.yml down
```

## Test 3: Measure payload sizes and compression (no prod traffic)

Sends batches to a local HTTP server that logs wire size, decompressed size, compression ratio, and decompression time. No data is sent to production.

Generates **realistic multi-service traffic** with all three signal types:
- **Logs**: 3 services (api-gateway, order-processor, auth-service) with varied bodies, severity levels, and attributes
- **Traces**: 2 services (api-gateway, payment-service) with different span durations and status codes
- **Metrics**: 1 service (infra-monitor) with Kubernetes attributes

Tests 4 compression formats side by side:
- `none` — raw protobuf (what raincatcher publishes to NATS today)
- `snappy` — proposed for raincatcher→NATS (E-1029)
- `zstd` — used by STEF dual-publish
- `gzip` — default OTel HTTP transport compression

### Large batch (100K)

```bash
docker compose -f docker-compose-large-batch-sizer.yml up
docker logs -f payload-sizer-large
```

### Small batch (5K)

```bash
docker compose -f docker-compose-small-batch-sizer.yml up
docker logs -f payload-sizer-small
```

### Example output

```
[  none] [   logs] wire:   27,100,000 bytes ( 25.84 MB) | decompressed:   27,100,000 bytes ( 25.84 MB) | ratio:   1.0x | decompress: N/A
[snappy] [   logs] wire:    1,928,682 bytes (  1.84 MB) | decompressed:   27,100,000 bytes ( 25.84 MB) | ratio:  14.1x | decompress: 1234us
[  zstd] [   logs] wire:      396,527 bytes (  0.38 MB) | decompressed:   27,100,000 bytes ( 25.84 MB) | ratio:  68.3x | decompress: 567us
[  gzip] [   logs] wire:      543,860 bytes (  0.52 MB) | decompressed:   27,100,000 bytes ( 25.84 MB) | ratio:  49.8x | decompress: 890us
[  none] [ traces] wire:    ...
[snappy] [ traces] wire:    ...
[  none] [metrics] wire:    ...
```

### How to interpret

The **decompressed size** = raw protobuf = what raincatcher publishes to NATS today.
The **wire size per format** = what the NATS message size would be if raincatcher compressed before publishing.

Compare against pipeline limits:

| Limit | Value |
|-------|-------|
| NATS `max_payload` (dev) | 1 MB |
| NATS `max_payload` (prod) | 8 MB |
| `og-logs` `max_msg_size` (prod) | 2 MiB |
| `og-traces` `max_msg_size` (prod) | 10 MiB |
| `og-metrics` `max_msg_size` (prod) | 8 MiB |

## Results (100K batch, realistic multi-service traffic)

Tested with multiple telemetrygen services producing diverse resource attributes, log bodies, span types, and metric types.

### Payload sizes per signal and compression format

| Signal | Raw protobuf | Snappy | Zstd | Gzip |
|--------|-------------|--------|------|------|
| **Logs** | 57.63 MB | 4.36 MB (13.2x) | 0.63 MB (91.5x) | 0.91 MB (63.4x) |
| **Traces** | 27.34 MB | 4.04 MB (6.8x) | 2.55 MB (10.7x) | 2.73 MB (10.0x) |
| **Metrics** | 40.05 MB | 3.10 MB (12.9x) | 0.77 MB (52.0x) | 1.02 MB (39.4x) |

### Decompression time (Python/cramjam, indicative only)

| Signal | Snappy | Zstd | Gzip |
|--------|--------|------|------|
| **Logs** | 25-44 ms | 22-44 ms | 18-21 ms |
| **Traces** | 15 ms | 17 ms | 29 ms |
| **Metrics** | 6-19 ms | 14-46 ms | 15-25 ms |

> Note: These are Python decompression times, not Go. Go's snappy/zstd implementations are significantly faster (typically 5-10x). These numbers are useful for relative comparison between formats, not absolute performance.

### Would compression fit within NATS limits?

| Signal | Raw protobuf | Snappy compressed | vs prod `max_payload` (8 MB) | vs prod stream `max_msg_size` |
|--------|-------------|-------------------|-----------------------------|-----------------------------|
| **Logs** | 57.63 MB | **4.36 MB** | OK | **over** (og-logs: 2 MiB) |
| **Traces** | 27.34 MB | **4.04 MB** | OK | OK (og-traces: 10 MiB) |
| **Metrics** | 40.05 MB | **3.10 MB** | OK | OK (og-metrics: 8 MiB) |

### Key findings

1. **Snappy compression makes 100K batches fit within NATS `max_payload` (8 MB)** for all signal types
2. **Logs still exceed `og-logs` stream limit** (2 MiB) even with Snappy — the stream limit must be raised or batch size reduced
3. **Traces compress the least** (6.8x vs 13.2x for logs) because trace IDs and span IDs are high-entropy random bytes
4. **Zstd achieves the best ratios** (10-91x) but is slower; Snappy is the best CPU/compression tradeoff
5. **Without compression**, a 100K batch produces 27-58 MB of raw protobuf — 3.5-7x over the 8 MB `max_payload`

## Configs

| File | Batch Size | Target | Purpose |
|------|-----------|--------|---------|
| `otel-collector-config-large-batch.yaml` | 100,000 | Production | Reproduce 400 errors |
| `otel-collector-config-small-batch.yaml` | 5,000 | Production | Verify fix |
| `otel-collector-config-large-batch-sizer.yaml` | 100,000 | Local sizer | Measure payload sizes (all signals, all compression) |
| `otel-collector-config-small-batch-sizer.yaml` | 5,000 | Local sizer | Measure payload sizes (all signals, all compression) |

## Collector version

Uses `0.131.0` to match the partner's version (`service.version: "0.131.0-dev"` from their error logs).

## Cleanup

```bash
docker compose -f docker-compose-large-batch.yml down -v
docker compose -f docker-compose-small-batch.yml down -v
docker compose -f docker-compose-large-batch-sizer.yml down -v
docker compose -f docker-compose-small-batch-sizer.yml down -v
```
