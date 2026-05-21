#!/usr/bin/env bash
set -euo pipefail

API_KEY="${OLLY_GARDEN_API_KEY:?ERROR: OLLY_GARDEN_API_KEY is not set}"
ENDPOINT="https://in.ollygarden.cloud/v1/logs"

echo "============================================"
echo " Test 1: Minimal valid OTLP log request"
echo " Expects: 200 OK"
echo "============================================"
echo ""

curl -sv -X POST "$ENDPOINT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"batch-debug-test"}}]},"scopeLogs":[{"scope":{"name":"test"},"logRecords":[{"timeUnixNano":"1710000000000000000","body":{"stringValue":"hello from debug test"},"severityText":"INFO"}]}]}]}' \
  2>&1

echo ""
echo ""
echo "============================================"
echo " Test 2: Large payload (~1MB of log records)"
echo " Tests server body size limit"
echo "============================================"
echo ""

# Generate a ~1MB JSON payload with many log records
python3 -c "
import json
records = []
for i in range(5000):
    records.append({
        'timeUnixNano': '1710000000000000000',
        'body': {'stringValue': f'large batch test log entry {i} with padding to increase size ' + 'x' * 100},
        'severityText': 'INFO'
    })
payload = {'resourceLogs': [{'resource': {'attributes': [{'key': 'service.name', 'value': {'stringValue': 'batch-debug-test'}}]}, 'scopeLogs': [{'scope': {'name': 'test'}, 'logRecords': records}]}]}
print(json.dumps(payload))
" > /tmp/large-batch-payload.json

PAYLOAD_SIZE=$(wc -c < /tmp/large-batch-payload.json | tr -d ' ')
echo "Payload size: ${PAYLOAD_SIZE} bytes"
echo ""

curl -sv -X POST "$ENDPOINT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d @/tmp/large-batch-payload.json \
  2>&1

echo ""
echo ""
echo "============================================"
echo " Test 3: Very large payload (~5MB)"
echo " Likely to trigger 400 or 413"
echo "============================================"
echo ""

python3 -c "
import json
records = []
for i in range(25000):
    records.append({
        'timeUnixNano': '1710000000000000000',
        'body': {'stringValue': f'very large batch test log entry {i} with extra padding ' + 'x' * 100},
        'severityText': 'INFO'
    })
payload = {'resourceLogs': [{'resource': {'attributes': [{'key': 'service.name', 'value': {'stringValue': 'batch-debug-test'}}]}, 'scopeLogs': [{'scope': {'name': 'test'}, 'logRecords': records}]}]}
print(json.dumps(payload))
" > /tmp/very-large-batch-payload.json

PAYLOAD_SIZE=$(wc -c < /tmp/very-large-batch-payload.json | tr -d ' ')
echo "Payload size: ${PAYLOAD_SIZE} bytes"
echo ""

curl -sv -X POST "$ENDPOINT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d @/tmp/very-large-batch-payload.json \
  2>&1

echo ""
echo ""
echo "============================================"
echo " Test 4: Empty body"
echo " Checks how server handles invalid input"
echo "============================================"
echo ""

curl -sv -X POST "$ENDPOINT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '' \
  2>&1

echo ""
echo ""
echo "============================================"
echo " Done. Check the response headers above:"
echo "  - 'Server: nginx' = nginx rejected it"
echo "  - gRPC/app error  = raincatcher rejected it"
echo "  - Look for 'content-length' limits"
echo "============================================"

rm -f /tmp/large-batch-payload.json /tmp/very-large-batch-payload.json
