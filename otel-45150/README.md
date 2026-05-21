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

## How It Works

1. **Prometheus** scrapes itself and adds `_generated=true` label to specific metrics
2. **Remote write** sends to OTel Collector with relabel config that keeps:
   - `target_info` (always)
   - Metrics where `_generated=true`
3. **Bug trigger**: When a batch contains only `target_info` (no other matching metrics), the collector responds with `0 accepted samples`
4. **Result**: Prometheus interprets this as a failure since it sent samples but none were accepted


# Check OTel Collector is receiving data
docker logs otel-collector


.yaml that I used as example: https://github.com/prometheus/prometheus/blob/main/config/testdata/conf.good.yml