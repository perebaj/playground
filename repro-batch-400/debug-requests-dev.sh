#!/usr/bin/env bash
set -euo pipefail

API_KEY="${OLLY_GARDEN_API_KEY:?ERROR: OLLY_GARDEN_API_KEY is not set}"
ENDPOINT="https://in-dev.ollygarden.cloud/v1/logs"

# Helper: generate a JSON OTLP log payload with N records, each with ~200 bytes of body
generate_payload() {
  local count=$1
  local label=$2
  python3 -c "
import json
records = []
for i in range($count):
    records.append({
        'timeUnixNano': '1710000000000000000',
        'body': {'stringValue': f'[$label] pipeline limit test entry {i} ' + 'x' * 150},
        'severityText': 'INFO',
        'attributes': [
            {'key': 'test.batch_size', 'value': {'intValue': str($count)}},
            {'key': 'test.label', 'value': {'stringValue': '$label'}}
        ]
    })
payload = {
    'resourceLogs': [{
        'resource': {
            'attributes': [
                {'key': 'service.name', 'value': {'stringValue': 'dev-limit-test'}},
                {'key': 'service.namespace', 'value': {'stringValue': 'e2e-dev-limit-test'}}
            ]
        },
        'scopeLogs': [{
            'scope': {'name': 'limit-test'},
            'logRecords': records
        }]
    }]
}
print(json.dumps(payload))
"
}

run_test() {
  local label=$1
  local count=$2
  local description=$3

  echo "============================================"
  echo " $label: $description"
  echo "============================================"

  local payload_file
  payload_file=$(mktemp)
  generate_payload "$count" "$label" > "$payload_file"

  local payload_size
  payload_size=$(wc -c < "$payload_file" | tr -d ' ')
  local payload_mb
  payload_mb=$(echo "scale=2; $payload_size / 1048576" | bc)
  echo "  JSON payload: ${payload_size} bytes (${payload_mb} MB)"

  # Send uncompressed
  echo "  Sending uncompressed..."
  local http_code
  http_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$ENDPOINT" \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    -d @"$payload_file" \
    --max-time 30 2>&1) || http_code="TIMEOUT"
  echo "  -> HTTP $http_code (uncompressed)"

  # Send gzip compressed
  local compressed_file
  compressed_file=$(mktemp)
  gzip -c "$payload_file" > "$compressed_file"
  local compressed_size
  compressed_size=$(wc -c < "$compressed_file" | tr -d ' ')
  local compressed_mb
  compressed_mb=$(echo "scale=2; $compressed_size / 1048576" | bc)
  echo "  Sending gzip compressed (${compressed_size} bytes / ${compressed_mb} MB)..."
  http_code=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$ENDPOINT" \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    -H "Content-Encoding: gzip" \
    --data-binary @"$compressed_file" \
    --max-time 30 2>&1) || http_code="TIMEOUT"
  echo "  -> HTTP $http_code (gzip)"

  rm -f "$payload_file" "$compressed_file"
  echo ""
}

echo ""
echo "Target: $ENDPOINT"
echo "Testing dev pipeline limits (nginx=64M, raincatcher max_body_size, NATS max_payload=64MB)"
echo ""

# Test 1: Small payload — should always work
run_test "T1-small" 100 "100 records (~25 KB) - baseline"

# Test 2: Medium payload — under old 1MB limit
run_test "T2-medium" 3000 "3K records (~750 KB) - under old 1MB limit"

# Test 3: Over old 1MB limit but under 10MB
run_test "T3-over-1mb" 5000 "5K records (~1.2 MB) - over old og-logs 1MB limit"

# Test 4: Over old 10MB limit but under 64MB
run_test "T4-over-10mb" 50000 "50K records (~12 MB) - over old max_payload 10MB"

# Test 5: Large batch — close to 64MB limit
run_test "T5-large" 200000 "200K records (~50 MB) - near 64MB limit"

# Test 6: Very large — should exceed 64MB and get rejected
run_test "T6-over-limit" 300000 "300K records (~75 MB) - should exceed 64MB limit"

echo "============================================"
echo " Summary"
echo "============================================"
echo " T1-T5: Expected 200 (within 64MB pipeline limit)"
echo " T6:    Expected 413 or 400 (exceeds 64MB limit)"
echo "============================================"
