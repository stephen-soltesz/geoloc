select
  latlon, index

FROM (
  SELECT
    ROW_NUMBER() OVER() as row,
    latlon,
    CASE
    WHEN download BETWEEN 0 AND 10 THEN 0
    WHEN download BETWEEN 10 AND 25 THEN 1
    WHEN download BETWEEN 25 AND 50 THEN 2
    WHEN download BETWEEN 50 AND 100 THEN 3
    WHEN download BETWEEN 100 AND 10000 THEN 4
    ELSE 5
  END as index

FROM (
SELECT
  metro_raw,
  APPROX_QUANTILES(download, 10)[OFFSET(5)] AS download,
  FORMAT("%f,%f", lat, lon) AS latlon
FROM (
  SELECT
    REGEXP_EXTRACT(connection_spec.server_hostname, r"mlab\d.([a-z]{3})..") AS metro_raw,
    (8 * (web100_log_entry.snap.HCThruOctetsAcked / (
        web100_log_entry.snap.SndLimTimeRwin +
        web100_log_entry.snap.SndLimTimeCwnd +
        web100_log_entry.snap.SndLimTimeSnd))) AS download,
    connection_spec.client_geolocation.latitude AS lat,
    connection_spec.client_geolocation.longitude AS lon
  FROM
    `measurement-lab.base_tables.ndt`
  WHERE
    _PARTITIONTIME >= TIMESTAMP_SUB(TIMESTAMP_TRUNC(CURRENT_TIMESTAMP(), DAY), INTERVAL 24 HOUR)
    AND connection_spec.data_direction = 1
    AND connection_spec.tls IS NOT FALSE
    AND web100_log_entry.snap.HCThruOctetsAcked >= 8192
    AND (web100_log_entry.snap.SndLimTimeRwin + web100_log_entry.snap.SndLimTimeCwnd + web100_log_entry.snap.SndLimTimeSnd) >= 9000000
    AND (web100_log_entry.snap.SndLimTimeRwin + web100_log_entry.snap.SndLimTimeCwnd + web100_log_entry.snap.SndLimTimeSnd) < 600000000
    AND connection_spec.client_geolocation.latitude IS NOT NULL
    AND connection_spec.client_geolocation.longitude IS NOT NULL
  GROUP BY
    metro_raw,
    download,
    lat,
    lon
)
GROUP BY metro_raw, latlon
)
)
WHERE MOD(row, 3) = 0
