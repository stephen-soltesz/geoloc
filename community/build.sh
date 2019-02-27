#!/bin/bash


set -x

gopherjs build worldmap.go
cat dscc.min.js sites.js worldmap.js worldMap.js > viz.js
gsutil -h 'Cache-Control:private, max-age=0, no-transform' cp -a public-read \
    viz.js viz.json viz.css manifest.json gs://soltesz-mlab-sandbox/v2/

