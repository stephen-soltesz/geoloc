# Precompution for browser clients

## Pipeline Steps

* Run query to download recent set of client locations.

  - [+] deduplicate locations.

* Parse resulting JSON:

  - [+] convert string floats
  - [+] partition files by metro
  - [+] write partitioned files to resources directory

* Parse sites config information.

  - [-] normalize information: metro, lat, lon.
  - [-] create a list of those values.
  - [-] filter bad sites & convert ambiguous sites (change query also).

* On double click, fetch, parse and plot client points. Cache values.

  - [+] individual metro file downloads.

## Viz

* [+] Load sites via json
* [+] Draw sites automatically / smaller sites
* [+] Clear all button
* [+] map mouse movement
* [+] eliminate some unnecessary complexity (static load).
* [-] continue to simplify static load logic.
* [-] voronoi should wrap around.