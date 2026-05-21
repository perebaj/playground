#!/usr/bin/env bash
# Live view of replicas, ScaledObject status, HPA, and JetStream consumer lag.
# Refreshes every 3 seconds.
set -u

while true; do
  clear
  printf '\033[1;33m=== %s ===\033[0m\n\n' "$(date +'%H:%M:%S')"

  echo "── Consumer pods ──"
  kubectl -n demo get pods -l app=slow-consumer --no-headers 2>/dev/null \
    | awk '{ printf "  %-32s %-9s %s\n", $1, $3, $5 }'
  echo

  echo "── ScaledObject ──"
  kubectl -n demo get scaledobject slow-consumer \
    -o 'custom-columns=READY:.status.conditions[?(@.type=="Ready")].status,ACTIVE:.status.conditions[?(@.type=="Active")].status,FALLBACK:.status.conditions[?(@.type=="Fallback")].status,MIN:.spec.minReplicaCount,MAX:.spec.maxReplicaCount' \
    --no-headers 2>/dev/null | sed 's/^/  /'
  echo

  echo "── HPA (managed by KEDA) ──"
  kubectl -n demo get hpa keda-hpa-slow-consumer \
    -o 'custom-columns=CURRENT:.status.currentReplicas,DESIRED:.status.desiredReplicas,MIN:.spec.minReplicas,MAX:.spec.maxReplicas' \
    --no-headers 2>/dev/null | sed 's/^/  /' || echo "  (HPA not created yet — KEDA generates it on demand)"
  echo

  echo "── JetStream consumer lag ──"
  kubectl -n demo exec deploy/nats-tools -- nats --server="$NATS_URL" consumer info test-stream slow-consumer 2>/dev/null \
    | awk '/Outstanding Acks|Pending Messages|Last Delivered|Redelivered|Waiting Pulls/ { print "  " $0 }' || echo "  (nats-tools not ready)"

  sleep 3
done
