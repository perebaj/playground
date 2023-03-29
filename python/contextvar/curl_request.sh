#!/bin/bash



# Número de requisições a serem feitas
num_requests=10

# Loop para gerar valores aleatórios e executar a requisição curl
for (( i=0; i<$num_requests; i++ )); do
  org_id="org_id_$i"
  trace_id="trace_id_$i"
  (time curl -X 'GET' \
    'http://127.0.0.1:8000/' \
    -H 'accept: application/json' \
    -H "organization-id: $org_id" \
    -H "trace-id: $trace_id" &)
done