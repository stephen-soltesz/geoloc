#!/usr/bin/env jsonnet -J . --string
local sites = import 'mlab-site-stats.json';
local metros = std.set(
  [
    site.site[0:3]
    for site in sites
  ]
);

local conv(s) = {
  metro: s.site[0:3],
  lat: '%s' % s.latitude,
  lon: '%s' % s.longitude,
};

local collapse(a, b) = (
  local last = a[std.length(a) - 1];
  local c = conv(b);
  if last.metro == c.metro then
    a
  else
    a + [c]
);

local unique = std.foldl(collapse, sites[1:], [conv(sites[0])]);

'sites = ' + std.toString(unique) + ';'
