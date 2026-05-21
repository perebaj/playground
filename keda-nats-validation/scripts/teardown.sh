#!/usr/bin/env bash
set -euo pipefail
CLUSTER_NAME=keda-nats-validation
kind delete cluster --name "$CLUSTER_NAME"
