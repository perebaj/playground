#!/bin/bash
# Boots a Quack server with per-org authentication: a small token→org table
# drives the custom authentication callback. We START with permissive authz
# and tighten in a second pass once we observe what arguments each callback
# actually receives.
set -euo pipefail
cd "$(dirname "$0")/.."

{
  cat <<'SQL'
LOAD quack;
LOAD httpfs;

-- Token → org mapping. The custom authentication callback consults this.
CREATE TABLE token_orgs(token VARCHAR PRIMARY KEY, allowed_org VARCHAR);
INSERT INTO token_orgs VALUES
  ('probe-alpha-token', 'org_test_alpha'),
  ('probe-beta-token',  'org_test_beta'),
  ('probe-demo-token',  'org_test_demo');

-- Authentication: 3-arg callback returning BOOLEAN. We don't yet know which
-- arg holds the token, so return true if ANY of the three args matches a
-- known token. The follow-up probe (m1-args-probe.sh) will narrow it down.
CREATE MACRO m1_authn(a, b, c) AS (
  EXISTS(SELECT 1 FROM token_orgs WHERE token IN (a, b, c))
);

-- Authorization: 2-arg callback returning BOOLEAN. Permissive for now —
-- the goal of the first server boot is to learn the callback signatures.
CREATE MACRO m1_authz(a, b) AS (true);

SET quack_authentication_function = 'm1_authn';
SET quack_authorization_function  = 'm1_authz';

CREATE OR REPLACE VIEW traces AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/traces/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW logs AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/logs/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_gauge AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_gauge/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_sum AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_sum/date=*/data.parquet', hive_partitioning=true);
CREATE OR REPLACE VIEW metrics_histogram AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/metrics_histogram/date=*/data.parquet', hive_partitioning=true);

-- We deliberately omit the `token :=` argument so the built-in token check is
-- defaulted; our custom m1_authn supersedes it via the authentication function.
CALL quack_serve('quack:localhost:9494', allow_other_hostname := true);
SQL
  tail -f /dev/null
} | duckdb
