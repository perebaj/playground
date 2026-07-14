#!/bin/bash
# Reproduce the QuackScan column-narrowing crash locally.
#
# The bug: a client query that returns fewer columns than the source
# catalog object (projection, count(*), aggregation) crashes the client
# with an `InternalException` in QuackScan — but ONLY when the source
# is a VIEW wrapping a table function like read_parquet. The identical
# data exposed as a materialized TABLE has no such threshold and
# returns correct results at every volume tested.
#
# Full write-up: M2-QUACKSCAN-BUG.md.
#
# Requires: docker, duckdb CLI (>=1.5) with `INSTALL quack` reachable.
# Run from the repo root:
#
#   scripts/repro-quackscan-bug.sh
#
# Exits non-zero if the pattern doesn't hold — either the baseline
# passes fail (small volume or full-schema queries stop working) or
# the trigger queries stop crashing (bug fixed upstream) or the
# TABLE-backed control queries crash (invalidates the workaround).

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

-- Same underlying data exposed both ways so we can prove the bug is
-- view-specific, not projection-specific.
CREATE OR REPLACE VIEW  v_small AS SELECT * FROM read_parquet('/data/small.parquet');
CREATE OR REPLACE VIEW  v_over  AS SELECT * FROM read_parquet('/data/over.parquet');
CREATE OR REPLACE TABLE t_over  AS SELECT * FROM read_parquet('/data/over.parquet');

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
  printf "    [%-56s] %s\n" "$label" "$st"
  [ "$st" = "ok" ]
}

BASELINE_OK=true
CRASH_HIT=false
TABLE_CONTROL_OK=true

echo "==> Baseline: below threshold, VIEW projection works"
run_q "small VIEW: SELECT c1"                             "SELECT c1 FROM q.v_small;"                       || BASELINE_OK=false
run_q "small VIEW: SELECT count(*)"                       "SELECT count(*) FROM q.v_small;"                 || BASELINE_OK=false

echo "==> Baseline: full-schema passthrough always works (even past threshold)"
run_q "over VIEW: SELECT *"                               "SELECT * FROM q.v_over;"                         || BASELINE_OK=false
run_q "over VIEW: SELECT c1, c2"                          "SELECT c1, c2 FROM q.v_over;"                    || BASELINE_OK=false
run_q "over VIEW: SELECT c1+c2 AS s (derived, ok)"        "SELECT c1+c2 AS s FROM q.v_over;"                || BASELINE_OK=false

echo "==> Control: same data as TABLE — projection is fine at any volume"
run_q "over TABLE: SELECT c1"                             "SELECT c1 FROM q.t_over;"                        || TABLE_CONTROL_OK=false
run_q "over TABLE: SELECT count(*)"                       "SELECT count(*) FROM q.t_over;"                  || TABLE_CONTROL_OK=false

echo "==> Trigger: VIEW-backed column narrowing past threshold — expect CRASH"
run_q "over VIEW: SELECT c1 (EXPECT CRASH)"               "SELECT c1 FROM q.v_over;"                        && CRASH_HIT=true || true
run_q "over VIEW: SELECT count(*) (EXPECT CRASH)"         "SELECT count(*) FROM q.v_over;"                  && CRASH_HIT=true || true

echo ""
if ! $BASELINE_OK; then
  echo "!! Baseline queries did not all succeed — either bug fixed upstream or repro broken."
  exit 2
fi
if ! $TABLE_CONTROL_OK; then
  echo "!! TABLE-backed control queries crashed — invalidates the view-specific hypothesis."
  exit 3
fi
if $CRASH_HIT; then
  echo "!! A VIEW-backed trigger query returned OK — bug may be fixed upstream."
  exit 4
fi

echo "OK — pattern reproduced:"
echo "   * VIEW-backed column narrowing past threshold crashes"
echo "   * TABLE-backed same-shape queries succeed"
echo "   * VIEW-backed full-schema passthrough succeeds"
