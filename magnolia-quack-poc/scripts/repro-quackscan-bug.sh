#!/bin/bash
# Reproduce the QuackScan column-narrowing crash locally.
#
# The bug: any query that asks the server to return a DataChunk with
# fewer columns than the underlying view (projection, aggregation)
# crashes the client with an `InternalException` in QuackScan once the
# underlying data crosses a row-count threshold. `SELECT * FROM view`
# (full-schema passthrough) works at every volume.
#
# Full write-up: M2-QUACKSCAN-BUG.md.
#
# Requires: docker, duckdb CLI (>=1.5) with `INSTALL quack` reachable.
# Run from the repo root:
#
#   scripts/repro-quackscan-bug.sh
#
# Exits non-zero if the pattern doesn't hold (either baseline queries
# fail, or the trigger queries stop crashing).

set -euo pipefail

cd "$(dirname "$0")/.."

TOKEN="localtoken"
IMAGE="${QUACK_IMAGE:-ghcr.io/ollygarden/magnolia/quack-duckdb:0.0.48}"
CONTAINER="quack-repro"

echo "==> Generating minimal fixtures (2 BIGINT columns, varying row counts)"
duckdb <<'SQL' >/dev/null
COPY (SELECT range::BIGINT AS c1, range::BIGINT AS c2 FROM range(22528))
  TO '/tmp/repro_small.parquet' (FORMAT PARQUET);
COPY (SELECT range::BIGINT AS c1, range::BIGINT AS c2 FROM range(26624))
  TO '/tmp/repro_over.parquet' (FORMAT PARQUET);
SQL

cat > /tmp/repro_init.sql << 'SQL'
INSTALL quack; INSTALL httpfs;
LOAD quack; LOAD httpfs;
CREATE OR REPLACE VIEW small AS SELECT * FROM read_parquet('/data/small.parquet');
CREATE OR REPLACE VIEW over  AS SELECT * FROM read_parquet('/data/over.parquet');
CALL quack_serve('quack:[::]:9494', token := 'localtoken', allow_other_hostname := true);
SQL

echo "==> Starting quack server in Docker"
docker rm -f "$CONTAINER" >/dev/null 2>&1 || true
docker run -d --name "$CONTAINER" --platform linux/amd64 \
  -v /tmp/repro_small.parquet:/data/small.parquet:ro \
  -v /tmp/repro_over.parquet:/data/over.parquet:ro \
  -v /tmp/repro_init.sql:/init/init.sql:ro \
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
  printf "    [%-52s] %s\n" "$label" "$st"
  [ "$st" = "ok" ]
}

BASELINE_OK=true
CRASH_HIT=false

echo "==> Baseline: below threshold, everything works"
run_q "small: SELECT c1 FROM q.small"                     "SELECT c1 FROM q.small;"                     || BASELINE_OK=false
run_q "small: SELECT count(*) FROM q.small"               "SELECT count(*) FROM q.small;"               || BASELINE_OK=false

echo "==> Baseline: full-schema passthrough works at every volume"
run_q "over: SELECT * FROM q.over (WORKS)"                "SELECT * FROM q.over;"                       || BASELINE_OK=false
run_q "over: SELECT c1, c2 FROM q.over (WORKS)"           "SELECT c1, c2 FROM q.over;"                  || BASELINE_OK=false
run_q "over: SELECT c1+c2 AS s FROM q.over (WORKS)"       "SELECT c1+c2 AS s FROM q.over;"              || BASELINE_OK=false

echo "==> Trigger: column narrowing past the threshold — expect CRASH"
run_q "over: SELECT c1 FROM q.over (EXPECT CRASH)"        "SELECT c1 FROM q.over;"                      && CRASH_HIT=true || true
run_q "over: SELECT count(*) FROM q.over (EXPECT CRASH)"  "SELECT count(*) FROM q.over;"                && CRASH_HIT=true || true

echo ""
if ! $BASELINE_OK; then
  echo "!! Baseline queries did not all succeed — either bug fixed upstream or repro broken."
  exit 2
fi
if $CRASH_HIT; then
  echo "!! A trigger query returned OK — bug may be fixed upstream."
  exit 3
fi

echo "OK — baseline passes, column-narrowing over threshold crashes as expected. Bug reproduced."
