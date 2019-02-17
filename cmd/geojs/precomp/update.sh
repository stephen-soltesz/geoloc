#!/bin/bash

set -e

QUERY=clients.sql
BQ_RESULT=clients2.json

echo "Running query:"
cat ${QUERY}
echo "======================"
cat ${QUERY} | \
    bq --project measurement-lab query --format=prettyjson \
       --nouse_legacy_sql --max_rows=4000000 > "${BQ_RESULT}"

echo "Parsing clients to metros"
go run clientinfo/parseClients.go \
  -input ${BQ_RESULT} \
  -outdir ../resources/metros

go run serverinfo/parseServers.go \
  -input "metros.json" \
  -output ../resources/sites.json

echo "ok"