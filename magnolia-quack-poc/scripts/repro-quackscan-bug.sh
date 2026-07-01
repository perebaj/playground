#!/bin/bash
# Reproduce the QuackScan large-aggregation crash locally.
#
# Bug: SELECT count(*) FROM (SELECT * FROM q.traces LIMIT N) crashes the
# client with an `InternalException` in QuackScan somewhere between
# N=20000 and N=50000. Full write-up: M2-QUACKSCAN-BUG.md
#
# Requires: docker, duckdb CLI (>=1.5) with quack extension installable.
# Run from the repo root:
#
#   scripts/repro-quackscan-bug.sh
#
# The script generates fixtures with 60k traces per parquet file (well
# past the crash threshold), starts the quack server in Docker, and runs
# a threshold sweep against it. Exits non-zero if either the small-N
# baseline crashes or the large-N run does NOT crash — i.e. if the bug
# stops reproducing.

set -euo pipefail

cd "$(dirname "$0")/.."

TOKEN="localtoken"
IMAGE="${QUACK_IMAGE:-ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48}"
CONTAINER="quack-repro"

echo "==> Building fixture-gen"
go build -o bin/fixture-gen ./cmd/fixture-gen/

echo "==> Regenerating fixtures with 60k traces per file (may take ~10s)"
rm -rf fixtures
./bin/fixture-gen -traces 60000 -out ./fixtures > /dev/null
echo "    total fixture size: $(du -sh fixtures | cut -f1)"

echo "==> Writing init.sql"
cat > /tmp/init-repro.sql << 'SQL'
INSTALL quack; INSTALL httpfs;
LOAD quack; LOAD httpfs;
CREATE OR REPLACE VIEW traces AS
  SELECT * FROM read_parquet('/fixtures/org_id=*/tables/traces/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_gauge AS
  SELECT * FROM read_parquet('/fixtures/org_id=*/tables/metrics_gauge/date=*/data.parquet', hive_partitioning=true);
CALL quack_serve('quack:[::]:9494', token := 'localtoken', allow_other_hostname := true);
SQL

echo "==> Starting quack server in Docker"
docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" --platform linux/amd64 \
  -v "$PWD/fixtures:/fixtures:ro" \
  -v /tmp/init-repro.sql:/init/init.sql:ro \
  -p 9494:9494 \
  --entrypoint /bin/sh \
  "$IMAGE" \
  -c 'set -eu; rm -f /tmp/duckdb.stdin; mkfifo /tmp/duckdb.stdin; tail -f /dev/null > /tmp/duckdb.stdin & HOME=/tmp exec /usr/local/bin/duckdb -init /init/init.sql < /tmp/duckdb.stdin' \
  >/dev/null
trap 'docker rm -f "$CONTAINER" >/dev/null 2>&1 || true' EXIT
sleep 6

if ! docker logs "$CONTAINER" 2>&1 | grep -q listen_url; then
  echo "!! Server didn't start. Logs:"
  docker logs "$CONTAINER" 2>&1 | tail -20
  exit 1
fi
echo "    server up, listening on :9494"

run_q() {
  local label="$1" sql="$2"
  local out
  out=$(duckdb -c "INSTALL quack; LOAD quack; ATTACH 'localhost:9494' AS q (TYPE quack, TOKEN '$TOKEN'); $sql" 2>&1)
  local st="ok"
  echo "$out" | grep -qE "InternalException|assertion" && st="CRASH"
  printf "    [%-48s] %s\n" "$label" "$st"
  # Return crash-as-status for the caller's assertions
  [ "$st" = "ok" ]
}

echo "==> Threshold sweep"
BASELINE_OK=true
CRASH_HIT=false

run_q "streaming ServiceName LIMIT 5"                     "SELECT ServiceName FROM q.traces LIMIT 5;"                                     || BASELINE_OK=false
run_q "count LIMIT 100"                                   "SELECT count(*) FROM (SELECT * FROM q.traces LIMIT 100);"                     || BASELINE_OK=false
run_q "count LIMIT 20000"                                 "SELECT count(*) FROM (SELECT * FROM q.traces LIMIT 20000);"                   || BASELINE_OK=false
run_q "count metrics_gauge (control)"                     "SELECT count(*) FROM q.metrics_gauge;"                                        || BASELINE_OK=false
run_q "count LIMIT 50000 (expect CRASH)"                  "SELECT count(*) FROM (SELECT * FROM q.traces LIMIT 50000);"                   && CRASH_HIT=true || true
run_q "count LIMIT 100000 (expect CRASH)"                 "SELECT count(*) FROM (SELECT * FROM q.traces LIMIT 100000);"                  && CRASH_HIT=true || true

echo ""
if ! $BASELINE_OK; then
  echo "!! Baseline queries (small N + control) did not all succeed — repro is broken."
  exit 2
fi
if $CRASH_HIT; then
  echo "!! At least one large-N query returned OK — bug may be fixed upstream."
  exit 3
fi

echo "OK — small-N passes, large-N crashes as expected. Bug reproduced."
