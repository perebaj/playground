.DEFAULT_GOAL := help

# ── Local POC: kind not required ─────────────────────────────────────────
# fixtures + a local Quack server + smoke queries. Kind manifests under
# k8s/ are next, but the local flow is sufficient to validate the model.

BIN := ./bin/fixture-gen
PID := /tmp/quack-server.pid
LOG := /tmp/quack-server.log

.PHONY: help
help: ## list available targets
	@awk 'BEGIN{FS=":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-22s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## build the Go fixture generator
	go build -o $(BIN) ./cmd/fixture-gen

.PHONY: fixtures
fixtures: build ## (re)generate the parquet fixtures under ./fixtures
	rm -rf fixtures
	$(BIN) --out ./fixtures

.PHONY: serve
serve: ## start the local Quack server in the background
	@./scripts/run-server.sh > $(LOG) 2>&1 & echo $$! > $(PID)
	@sleep 2
	@if lsof -i :9494 >/dev/null 2>&1; then \
	    echo "quack server up on :9494 (pid $$(cat $(PID))) — log: $(LOG)"; \
	else \
	    echo "quack server FAILED to bind :9494 — see $(LOG)"; cat $(LOG); exit 1; \
	fi

.PHONY: stop
stop: ## stop the local Quack server
	@if [ -f $(PID) ]; then pkill -f "run-server\|tail -f /dev/null" 2>/dev/null; pkill -f duckdb 2>/dev/null; rm -f $(PID); echo "stopped"; else echo "no pid file"; fi

.PHONY: smoke
smoke: ## run the smoke query suite against the running server
	./scripts/smoke.sh

.PHONY: all
all: fixtures serve smoke ## full local flow: build → fixtures → serve → smoke
	$(MAKE) stop

.PHONY: clean
clean: stop ## remove fixtures + binary
	rm -rf fixtures $(BIN) $(LOG) $(PID)

# ── Kind cluster flow (work in progress) ─────────────────────────────────

KIND_CLUSTER ?= magnolia-quack

.PHONY: kind-up
kind-up: ## boot the kind cluster
	kind create cluster --name $(KIND_CLUSTER) --config kind/cluster.yaml

.PHONY: kind-down
kind-down: ## tear down the kind cluster
	kind delete cluster --name $(KIND_CLUSTER)

.PHONY: kind-deploy
kind-deploy: ## load fixtures into the cluster and deploy Quack
	kubectl --context kind-$(KIND_CLUSTER) apply -f k8s/quack-server.yaml

.PHONY: kind-portforward
kind-portforward: ## expose Quack 9494 on localhost (foreground)
	kubectl --context kind-$(KIND_CLUSTER) port-forward svc/quack-server 9494:9494
