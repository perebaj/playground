# KEDA × NATS JetStream — local validation

Sandbox for exercising KEDA's `nats-jetstream` scaler end-to-end against a
real NATS JetStream consumer. Runs entirely on `kind`, no cloud.

## What it sets up

- 1-node `kind` cluster
- NATS server (JetStream enabled) in namespace `nats`
- KEDA operator in namespace `keda`
- `nats-tools` helper pod in `demo` (for running `nats` CLI via `kubectl exec`)
- JetStream stream `test-stream` (subjects `test.>`) + durable pull consumer `slow-consumer`
- Deployment `slow-consumer` (image `natsio/nats-box`) — pulls 1 message, acks, sleeps 5s
- `ScaledObject` targeting `slow-consumer` with `lagThreshold=10`, `min=1`, `max=5`

The combination — `lagThreshold=10`, each pod ~0.1-0.2 msg/s — makes scaling
visible within ~1 minute of publishing a small burst.

## Prerequisites

- `kind` 0.20+
- `kubectl` 1.28+
- `helm` 3.x
- `docker` (kind backend)

## Quick start

```bash
./scripts/setup.sh        # creates kind cluster + installs everything
./scripts/watch.sh        # (new terminal) live view of pods, lag, replicas
./scripts/publish.sh 100  # publish 100 messages to drive backlog
# watch the replicas climb in the other terminal
./scripts/publish.sh 0    # (no-op) wait, then observe scale-down after cooldown
./scripts/teardown.sh     # remove the kind cluster
```

## What to watch for

1. **Idle state** — 1 replica, lag near 0.
2. **After `publish.sh 100`** — within `pollingInterval` (10s), `num_pending` jumps
   to ~100, KEDA computes `desired = ceil(100 / 10) = 10`, clamped to `max=5`.
   HPA `behavior.scaleUp` allows +100% per 15s → replicas climb 1 → 2 → 4 → 5.
3. **As pods consume** — `num_pending` drops, `desired` falls.
4. **After backlog clears** — `desired = 1`. HPA's `scaleDown.stabilizationWindowSeconds=60`
   means HPA waits a minute of consistently low lag before each scale-down step.
   `cooldownPeriod=60` (KEDA) adds an extra dwell before the FIRST scale-down.
5. **Steady state** — back to 1 replica.

## Tuning knobs

- `manifests/04-scaledobject.yaml` — change `lagThreshold`, `min/maxReplicaCount`,
  `pollingInterval`, `cooldownPeriod`, `behavior.*` to experience different curves.
- `manifests/03-consumer.yaml` — increase the `sleep` to make consumers slower
  (longer backlog drains, more replicas).
- `scripts/publish.sh <N>` — publish bursts of any size.

## Inspection cheat sheet

All `nats` commands run via `kubectl exec` into the `nats-box` pod that the
NATS chart deploys in the `nats` namespace. Set this once per shell to
shorten the commands:

```bash
alias n='kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222'
```

(Examples below use the full form so they work without the alias.)

### NATS server health

```bash
# Connection + RTT + JetStream availability
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 server check connection
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 server check jetstream

# Uptime, memory, connection count, msgs in/out
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 server info

# All servers in the cluster (with N nodes you'd see leaders/peers here)
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 server ls

# JetStream-wide report (memory/disk per server, leaders, etc.)
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 server report jetstream
```

### Streams

```bash
# All streams: messages, bytes, lost, consumers, replicas
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 stream report

# One stream's full config + state
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 stream info test-stream

# Subject distribution inside the stream
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 stream subjects test-stream
```

### Consumers (the lag KEDA scales on)

`Unprocessed Messages` is the `num_pending` that KEDA reads via the
monitoring endpoint.

```bash
# Ack pending / unprocessed / redelivered for every consumer in the stream
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 consumer report test-stream

# One consumer: config + Outstanding Acks + Unprocessed + Last Delivered
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 consumer info test-stream slow-consumer

# Live view of lag draining
watch -n2 'kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 consumer info test-stream slow-consumer 2>/dev/null | grep -E "Outstanding|Unprocessed|Redelivered|Last Delivered"'
```

### Scaling state (KEDA + HPA)

```bash
# ScaledObject overview (READY / ACTIVE / triggers / min-max)
kubectl -n demo get scaledobject slow-consumer

# Per-trigger last error + ack metric value
kubectl -n demo describe scaledobject slow-consumer

# HPA generated by KEDA — current vs desired replicas + target metric
kubectl -n demo get hpa keda-hpa-slow-consumer
kubectl -n demo describe hpa keda-hpa-slow-consumer | tail -20

# Current pods of the scaled Deployment
kubectl -n demo get pods -l app=slow-consumer -o wide
```

### Operator logs (when scaling misbehaves)

```bash
# KEDA operator (scaling decisions, errors fetching metrics)
kubectl -n keda logs deploy/keda-operator -f

# KEDA metric-server (serves HPA's metric queries; errors here = HPA gets `<unknown>`)
kubectl -n keda logs deploy/keda-operator-metrics-apiserver -f
```

### Raw Prometheus metrics

For ad-hoc inspection of what KEDA exposes (useful when building the
Grafana dashboard):

```bash
kubectl -n keda port-forward deploy/keda-operator 8080:8080 &
curl -s localhost:8080/metrics | grep -E 'keda_scaler_(active|metrics_value|metrics_latency|detail_errors)'
```

### One-shot snapshot

Lag + replicas + ScaledObject status in a single print:

```bash
kubectl -n nats exec deploy/nats-box -- nats --server=nats://nats.nats.svc.cluster.local:4222 consumer report test-stream && \
kubectl -n demo get scaledobject,hpa,pods -l app=slow-consumer
```

## Troubleshooting

If `setup.sh` fails on the `nats-stream-setup` Job, check `kubectl -n demo
describe job nats-stream-setup` for the failure reason; the most common
cause is a `nats` CLI flag mismatch — the manifest uses `--defaults` plus
overrides only for the few fields that matter.

If the HPA shows `TARGETS: <unknown>/10`, KEDA hasn't done its first poll
yet (wait `pollingInterval` seconds) or it can't reach the monitoring
endpoint. The endpoint is `nats-headless.nats.svc.cluster.local:8222` —
the `nats` Service only exposes the client port 4222.
