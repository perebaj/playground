#!/bin/bash
# Connect to the local Quack server (run-server.sh keeps it alive at
# localhost:9494) and exercise the query patterns Olive will issue.
# Failure = non-zero exit.
set -euo pipefail
cd "$(dirname "$0")/.."

QUACK_URI="quack:localhost:9494"

run() {
  local label="$1" sql="$2"
  echo "── $label"
  duckdb -c "
    LOAD quack;
    CREATE OR REPLACE SECRET smoke (TYPE quack, TOKEN 'spike_token');
    SELECT * FROM quack_query('$QUACK_URI', \$\$ $sql \$\$);
  " 2>&1 | tail -n +5
  echo
}

run "1. table counts (5 views)" \
  "SELECT 'traces' t, count(*) n FROM traces
   UNION ALL SELECT 'logs', count(*) FROM logs
   UNION ALL SELECT 'metrics_gauge', count(*) FROM metrics_gauge
   UNION ALL SELECT 'metrics_sum', count(*) FROM metrics_sum
   UNION ALL SELECT 'metrics_histogram', count(*) FROM metrics_histogram
   ORDER BY 1;"

run "2. per-org-per-date span count (hive partitioning visible)" \
  "SELECT org_id, date, count(*) AS spans
   FROM traces GROUP BY 1,2 ORDER BY 1,2;"

run "3. typed endpoint /services for org_test_alpha" \
  "SELECT ServiceName, count(*) AS spans
   FROM traces
   WHERE org_id = 'org_test_alpha'
   GROUP BY 1 ORDER BY 2 DESC;"

run "4. typed /attributes/cloud.platform/distribution" \
  "SELECT
     json_extract_string(ResourceAttributes, '\$.\"cloud.platform\"') AS platform,
     count(*) AS spans
   FROM traces
   WHERE org_id = 'org_test_beta'
   GROUP BY 1 ORDER BY 2 DESC;"

run "5. partition pruning sanity: date filter narrows to one parquet" \
  "SELECT count(*) AS traces_2026_06_29
   FROM traces
   WHERE org_id = 'org_test_demo' AND date = '2026-06-29';"

run "6. cross-org admin query (NO org filter)" \
  "SELECT org_id, count(DISTINCT ServiceName) AS unique_services
   FROM traces GROUP BY 1 ORDER BY 2 DESC;"

run "7. per-org isolation check — simulate bad caller trying to cross orgs.
    The Olive layer is what enforces this; on the server, it executes." \
  "SELECT org_id, count(*) AS leaked_rows
   FROM traces
   WHERE org_id IN ('org_test_alpha', 'org_test_beta')
   GROUP BY 1;"

echo "── smoke complete"
