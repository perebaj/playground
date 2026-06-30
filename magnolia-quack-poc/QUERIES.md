# Hands-on queries against the Quack POC

Each block is a complete shell command — copy, paste, run. No context needed between blocks. The fixtures are deterministic (seeded by org_id + date), so numbers should be reproducible across re-runs.

## Prerequisites

You need one of:

- **Local mode**: `make serve` (Quack on `localhost:9494`, no kind needed)
- **Kind mode**: `make kind-up && make kind-deploy`, then `make kind-portforward` in another shell so `localhost:9494` reaches the pod.

Confirm the server is reachable:

```bash
lsof -i :9494
```

## 1. Schema discovery

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$ SHOW TABLES \$\$);
"
```

→ Five views: `logs`, `metrics_gauge`, `metrics_histogram`, `metrics_sum`, `traces`.

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$ DESCRIBE traces \$\$);
"
```

→ The trace schema (column, type, null/key/default). Note `date` and `org_id` at the bottom — those are hive partition keys, not stored inside the parquet files.

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT 'traces' t, count(*) c FROM traces UNION ALL
  SELECT 'logs', count(*) FROM logs UNION ALL
  SELECT 'metrics_gauge', count(*) FROM metrics_gauge UNION ALL
  SELECT 'metrics_sum', count(*) FROM metrics_sum UNION ALL
  SELECT 'metrics_histogram', count(*) FROM metrics_histogram
\$\$);
"
```

→ Total per signal: 300 / 600 / 1200 / 1200 / 300.

## 2. Per-org analytics

These mirror what the Olive frontend / Rose's report-magnolia-generate skill would issue.

### Top services in one org

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT ServiceName, count(*) AS spans
  FROM traces
  WHERE org_id = 'org_test_alpha'
  GROUP BY 1 ORDER BY 2 DESC
\$\$);
"
```

→ `checkout`, `cart`, `frontend` with counts (varies by random seed).

### Top scope names (vendor / instrumentation detection)

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT ScopeName, count(*) AS spans
  FROM traces
  WHERE org_id = 'org_test_beta'
  GROUP BY 1 ORDER BY 2 DESC
\$\$);
"
```

→ Two scopes (`go.opentelemetry.io/otel/sdk/tracer`, `github.com/XSAM/otelsql`).

### Cloud platform distribution per org

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    org_id,
    json_extract_string(ResourceAttributes, '\$.\"cloud.platform\"') AS platform,
    count(*) AS spans
  FROM traces
  GROUP BY 1, 2 ORDER BY 1, 3 DESC
\$\$);
"
```

→ alpha + demo on `gcp_kubernetes_engine`, beta on `aws_ec2`.

### k8s cluster per org (NULL for AWS-only)

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    org_id,
    json_extract_string(ResourceAttributes, '\$.\"k8s.cluster.name\"') AS cluster,
    count(*) AS spans
  FROM traces
  GROUP BY 1, 2 ORDER BY 1
\$\$);
"
```

→ alpha → `alpha-prod`, demo → `demo-cluster`, beta → NULL.

## 3. Multi-week (cross-date) analytics

### Service mix per (org, date)

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT org_id, date, ServiceName, count(*) AS spans
  FROM traces
  GROUP BY 1, 2, 3 ORDER BY 1, 2, 4 DESC
\$\$);
"
```

### Service count delta between week 1 and week 2

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  WITH per_week AS (
    SELECT date, ServiceName, count(*) AS spans
    FROM traces
    WHERE org_id = 'org_test_demo'
    GROUP BY 1, 2
  )
  SELECT
    ServiceName,
    sum(CASE WHEN date = '2026-06-22' THEN spans ELSE 0 END) AS w1,
    sum(CASE WHEN date = '2026-06-29' THEN spans ELSE 0 END) AS w2
  FROM per_week
  GROUP BY 1 ORDER BY 1
\$\$);
"
```

→ Three rows (ad, flagd, product-catalog).

## 4. Partition pruning evidence

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  EXPLAIN ANALYZE
  SELECT count(*) FROM traces
  WHERE org_id = 'org_test_demo' AND date = '2026-06-29'
\$\$);
"
```

→ The `READ_PARQUET` operator reports a single file scanned, not six.

## 5. JSON attribute drilldowns

### Service-namespace breakdown

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    json_extract_string(ResourceAttributes, '\$.\"k8s.namespace.name\"') AS ns,
    count(*) AS spans
  FROM traces
  WHERE org_id = 'org_test_alpha'
  GROUP BY 1 ORDER BY 2 DESC
\$\$);
"
```

→ Single namespace `default` for alpha.

### Span attribute drilldown (http.method, http.status_code)

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    json_extract_string(SpanAttributes, '\$.\"http.method\"') AS method,
    json_extract_string(SpanAttributes, '\$.\"http.status_code\"') AS status,
    count(*) AS spans
  FROM traces
  WHERE org_id = 'org_test_alpha'
  GROUP BY 1, 2 ORDER BY 3 DESC
\$\$);
"
```

→ One row: `GET / 200 / 100` (fixtures only generate one combo).

## 6. OpenTelemetry maturity signals

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    org_id,
    count(*) FILTER (WHERE json_extract_string(ResourceAttributes, '\$.\"telemetry.sdk.name\"') IS NOT NULL) AS with_sdk,
    count(*) FILTER (WHERE json_extract_string(ResourceAttributes, '\$.\"telemetry.sdk.name\"') IS NULL) AS without_sdk
  FROM traces
  GROUP BY 1 ORDER BY 1
\$\$);
"
```

→ Only `org_test_demo` has telemetry.sdk.name populated (matches the real OTel demo's profile).

## 7. Cross-signal cross-reference

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  WITH t AS (SELECT org_id, date, count(*) c FROM traces GROUP BY 1,2),
       l AS (SELECT org_id, date, count(*) c FROM logs GROUP BY 1,2),
       g AS (SELECT org_id, date, count(*) c FROM metrics_gauge GROUP BY 1,2)
  SELECT t.org_id, t.date, t.c AS traces, l.c AS logs, g.c AS gauges
  FROM t
  JOIN l ON t.org_id = l.org_id AND t.date = l.date
  JOIN g ON t.org_id = g.org_id AND t.date = g.date
  ORDER BY 1, 2
\$\$);
"
```

→ Six rows: 50 / 100 / 200 per (org, date).

## 8. Admin / cross-org queries

These would be gated behind admin auth in production.

### Org inventory

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    org_id,
    count(DISTINCT date) AS dates_seen,
    count(DISTINCT ServiceName) AS services,
    count(*) AS spans
  FROM traces
  GROUP BY 1 ORDER BY 4 DESC
\$\$);
"
```

→ All three test orgs, 2 dates each, 3 services each.

### Cross-org cloud distribution

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT
    json_extract_string(ResourceAttributes, '\$.\"cloud.platform\"') AS platform,
    count(DISTINCT org_id) AS orgs,
    count(*) AS spans
  FROM traces
  GROUP BY 1 ORDER BY 3 DESC
\$\$);
"
```

→ Two rows: `gcp_kubernetes_engine` (2 orgs, 200 spans) and `aws_ec2` (1 org, 100 spans).

## 9. Per-org isolation: the deliberate gap

The Quack server does not gate cross-org access on its own. Anything the client SQL says is what runs. Convince yourself:

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT org_id, count(*) FROM traces GROUP BY 1
\$\$);
"
```

→ Three rows. Server happily exposed all three orgs.

This is the design property covered in the project doc: per-org isolation has to be enforced by Olive's translation layer (rewriting WHERE on typed endpoints, gating SQL passthrough to admin-only). The next experiment (M1: per-org logical databases backed by per-org tokens) targets pushing the boundary into Quack itself.

## 10. Useful debug helpers

### Server-side settings

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT name, value FROM duckdb_settings() WHERE name LIKE 'quack%'
\$\$);
"
```

### Auth check function in use

```bash
duckdb -c "
LOAD quack;
CREATE OR REPLACE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$
  SELECT current_setting('quack_authentication_function')
\$\$);
"
```

### Clear client-side cache

```bash
duckdb -c "LOAD quack; SELECT quack_clear_cache();"
```

## 11. From inside the pod (kind mode only)

If you want to skip port-forward and talk to the server from inside the cluster:

```bash
POD=$(kubectl --context kind-magnolia-quack get pods -l app=quack-server -o name | head -1)
kubectl --context kind-magnolia-quack exec -it $POD -- /usr/local/bin/duckdb -c "
LOAD quack;
CREATE SECRET q (TYPE quack, TOKEN 'spike_token');
SELECT * FROM quack_query('quack:localhost:9494', \$\$ SELECT count(*) FROM traces \$\$);
"
```

## 12. Helper script (less typing)

The same boilerplate wrapped as `./scripts/q "<SQL>"`:

```bash
./scripts/q "SHOW TABLES"
./scripts/q "DESCRIBE traces"
./scripts/q "SELECT count(*) FROM traces"
./scripts/q "SELECT org_id, count(*) FROM traces GROUP BY 1 ORDER BY 1"
```

For multi-line queries:

```bash
./scripts/q "$(cat <<'SQL'
WITH per_week AS (
  SELECT date, ServiceName, count(*) AS spans
  FROM traces WHERE org_id = 'org_test_demo'
  GROUP BY 1, 2
)
SELECT ServiceName,
  sum(CASE WHEN date = '2026-06-22' THEN spans ELSE 0 END) AS w1,
  sum(CASE WHEN date = '2026-06-29' THEN spans ELSE 0 END) AS w2
FROM per_week GROUP BY 1 ORDER BY 1
SQL
)"
```

## 13. Cleanup

```bash
make stop          # local mode
make kind-down     # kind mode
```
