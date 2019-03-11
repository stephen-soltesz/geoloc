local sites = import 'sites.json';
/*
local dedup = {
    [site.metro]: {
        lat: site.lat,
        lon: site.lon,
    }
    for site in sites
};
*/

std.lines(std.set([
  'SELECT true AS anchor, "%s" as metro, "%s,%s" as latlon UNION ALL' % [site.metro, site.lat, site.lon]
  for site in sites[0:std.length(sites)-2]
]) + [
  local s = sites[std.length(sites)-1];
  'SELECT true AS anchor, "%s" as metro, "%s,%s" as latlon' % [s.metro, s.lat, s.lon]
])
