#!/bin/bash
# Probes which position (a, b, or c) of the 3-arg quack_authentication_function
# carries the client's token.
#
# Strategy: restart the server three times. Each boot accepts ONLY when one
# specific arg slot matches a known token. The client tries the same token
# each time and we see which boot accepts.
#
# Usage: ./m1-probe-arg.sh <a|b|c>
set -euo pipefail
cd "$(dirname "$0")/.."

SLOT=${1:?"usage: $0 <a|b|c>"}

case "$SLOT" in
  a) BODY="(a IN (SELECT token FROM token_orgs))" ;;
  b) BODY="(b IN (SELECT token FROM token_orgs))" ;;
  c) BODY="(c IN (SELECT token FROM token_orgs))" ;;
  *) echo "slot must be a, b, or c"; exit 1 ;;
esac

{
  cat <<SQL
LOAD quack;
LOAD httpfs;
CREATE TABLE token_orgs(token VARCHAR, allowed_org VARCHAR);
INSERT INTO token_orgs VALUES
  ('probe-alpha-token', 'org_test_alpha'),
  ('probe-beta-token',  'org_test_beta'),
  ('probe-demo-token',  'org_test_demo');
CREATE MACRO m1_authn(a, b, c) AS $BODY;
CREATE MACRO m1_authz(a, b) AS (true);
SET quack_authentication_function = 'm1_authn';
SET quack_authorization_function  = 'm1_authz';
CREATE OR REPLACE VIEW traces AS SELECT * FROM read_parquet('fixtures/org_id=*/tables/traces/date=*/data.parquet', hive_partitioning=true);
CALL quack_serve('quack:localhost:9494', allow_other_hostname := true);
SQL
  tail -f /dev/null
} | duckdb
