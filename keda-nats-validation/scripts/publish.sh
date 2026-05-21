#!/usr/bin/env bash
# Publish N messages to test.bench (no-op if N == 0).
# Usage: ./scripts/publish.sh [N]
set -euo pipefail

N="${1:-100}"

if [ "$N" -le 0 ]; then
  echo "N=$N — nothing to publish. (Wait and observe scale-down behaviour.)"
  exit 0
fi

echo "Publishing $N messages to test.bench via nats-tools pod..."
kubectl -n demo exec deploy/nats-tools -- sh -c "
  for i in \$(seq 1 $N); do
    nats --server=\$NATS_URL pub test.bench \"payload-\$i\" >/dev/null
  done
  echo
  echo 'Stream info:'
  nats --server=\$NATS_URL stream info test-stream | grep -E 'Messages|Bytes|First Sequence|Last Sequence' || true
  echo
  echo 'Consumer info:'
  nats --server=\$NATS_URL consumer info test-stream slow-consumer | grep -E 'Outstanding Acks|Pending Messages|Redelivered|Acknowledgement' || true
"
