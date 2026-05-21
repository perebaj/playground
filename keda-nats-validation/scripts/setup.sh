#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME=keda-nats-validation
HERE="$(cd "$(dirname "$0")/.." && pwd)"

step() { printf '\n\033[1;34m==> %s\033[0m\n' "$*"; }

step "Creating kind cluster '$CLUSTER_NAME'"
if kind get clusters 2>/dev/null | grep -qx "$CLUSTER_NAME"; then
  echo "(cluster already exists, skipping create)"
else
  kind create cluster --config "$HERE/kind-config.yaml"
fi

kubectl config use-context "kind-${CLUSTER_NAME}"

step "Applying namespaces"
kubectl apply -f "$HERE/manifests/00-namespaces.yaml"

step "Adding helm repos"
helm repo add nats https://nats-io.github.io/k8s/helm/charts/ >/dev/null 2>&1 || true
helm repo add kedacore https://kedacore.github.io/charts >/dev/null 2>&1 || true
helm repo update >/dev/null

step "Installing NATS (JetStream enabled)"
helm upgrade --install nats nats/nats \
  --namespace nats \
  --version '1.3.3' \
  --values "$HERE/manifests/01-nats-values.yaml" \
  --wait --timeout 5m

step "Installing KEDA"
helm upgrade --install keda kedacore/keda \
  --namespace keda \
  --version '2.18.0' \
  --wait --timeout 5m

step "Waiting for KEDA CRDs to be ready"
kubectl wait --for=condition=Established crd/scaledobjects.keda.sh --timeout=60s

step "Creating JetStream stream + consumer"
kubectl apply -f "$HERE/manifests/02-stream-setup-job.yaml"
kubectl -n demo wait --for=condition=complete job/nats-stream-setup --timeout=180s

step "Deploying slow-consumer + nats-tools"
kubectl apply -f "$HERE/manifests/03-consumer.yaml"
kubectl -n demo rollout status deploy/slow-consumer --timeout=120s
kubectl -n demo rollout status deploy/nats-tools --timeout=120s

step "Creating ScaledObject"
kubectl apply -f "$HERE/manifests/04-scaledobject.yaml"
sleep 3
kubectl -n demo get scaledobject slow-consumer

step "Ready. Run ./scripts/watch.sh in another terminal, then ./scripts/publish.sh 100"
