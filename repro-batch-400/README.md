# Batch Size 400 Error Reproduction

Reproduces 400 errors caused by a partner's OTel Collector config that batches 100K log records into a single HTTP export request.

## Hypothesis

The `batch/logs_export` processor with `send_batch_size: 100000` creates OTLP HTTP requests too large for raincatcher to accept, resulting in HTTP 400 (InvalidArgument).

## Setup

```bash
cd repro-batch-400
export OLLY_GARDEN_API_KEY="og_sk_..."
```

## Test 1: Large batch (100K) — should reproduce the 400 errors

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

## Configs

| File | Batch Size | Expected Result |
|------|-----------|-----------------|
| `otel-collector-config-large-batch.yaml` | 100,000 (double batch) | 400 errors |
| `otel-collector-config-small-batch.yaml` | 5,000 (single batch) | Success |

## What the log generator does

- 10 workers sending 5,000 logs/sec each = ~50K logs/sec total
- Runs for 120 seconds
- The collector batches these and exports as HTTP POST to `https://in.ollygarden.cloud/v1/logs`

## Collector version

Uses `0.131.0` to match the partner's version (`service.version: "0.131.0-dev"` from their error logs).

## Cleanup

```bash
docker compose -f docker-compose-large-batch.yml down -v
docker compose -f docker-compose-small-batch.yml down -v
```
