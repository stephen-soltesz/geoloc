#!/bin/bash

set -xe

mkdir -p build
cat vendor/dscc.min.js \
    vendor/d3.min.js \
    vendor/jquery-3.3.1.min.js \
    vendor/data.world.js \
    map.js > build/viz.js

gsutil -h 'Cache-Control:private, max-age=0, no-transform' cp -a public-read \
    build/viz.js viz.json viz.css manifest.json gs://soltesz-mlab-sandbox/wm1/

