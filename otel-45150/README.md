# OTEL-163 Reproduction

Local reproduction setup for [opentelemetry-collector-contrib#45150](https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/45150).

## The Bug

The `prometheusremotewrite` receiver incorrectly returns `0 accepted samples` when Prometheus sends `target_info` metrics alongside filtered data, causing Prometheus to report non-recoverable errors.

## Quick Start

```bash
cd repro-otel-163

# Start the stack
docker compose up

# Watch for errors in Prometheus logs (in another terminal)
docker logs -f prometheus 2>&1 | grep -i "error\|fail"
```

## Expected Error

After ~30 seconds, you should see errors like:

```
non-recoverable error
failedSampleCount=4
"sent v2 request with 4 samples... but PRW 2.0 response header statistics indicate 0 samples were accepted"
```

## How It Works

1. **Prometheus** scrapes itself and adds `_generated=true` label to specific metrics
2. **Remote write** sends to OTel Collector with relabel config that keeps:
   - `target_info` (always)
   - Metrics where `_generated=true`
3. **Bug trigger**: When a batch contains only `target_info` (no other matching metrics), the collector responds with `0 accepted samples`
4. **Result**: Prometheus interprets this as a failure since it sent samples but none were accepted

## Verify the Bug

```bash
# Check Prometheus remote write status
curl -s http://localhost:9091/api/v1/status/runtimeinfo | jq '.data.storageRetention'

# Check OTel Collector is receiving data
docker logs otel-collector 2>&1 | grep -i "metric"
```

## Workaround

Remove `target_info` from the relabel config (loses resource attributes):

```yaml
write_relabel_configs:
  - source_labels: [_generated]
    regex: "true"
    action: keep
```

## Cleanup

```bash
docker compose down -v
```

YAML example: https://github.com/prometheus/prometheus/blob/main/config/testdata/conf.good.yml