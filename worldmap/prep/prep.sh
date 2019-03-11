#!/bin/bash

set -ex
./metros.jsonnet > ../build/data.sites.js

cat clients.sql | bq --project measurement-lab query --format=prettyjson \
    --nouse_legacy_sql --max_rows=4000000 > clients.json
time jsonnet -J . --string dsccFormat.jsonnet > ../build/data.clients.js
