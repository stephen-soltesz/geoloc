SELECT
  row,
  lat,
  lon,
  metro
FROM (
  SELECT
    ROW_NUMBER() OVER(PARTITION BY metro) as row,
    c,
    m,
    CEIL(30000 * m/c) as lim,
    metro,
    lat,
    lon
  FROM (
    SELECT
        c,
        count(*) over(partition by metro_raw) as m,
        CASE WHEN metro_raw = "sjc" THEN "nuq" ELSE metro_raw END as metro,
        lat,
        lon
    FROM (
        SELECT
            count(*) over() as c,
            REGEXP_EXTRACT(connection_spec.server_hostname, r"mlab\d.([a-z]{3})..") AS metro_raw,
            connection_spec.client_geolocation.latitude AS lat,
            connection_spec.client_geolocation.longitude AS lon
        FROM
            `measurement-lab.base_tables.ndt`
        WHERE
            _PARTITIONTIME = TIMESTAMP_SUB(TIMESTAMP_TRUNC(CURRENT_TIMESTAMP(), DAY), INTERVAL 24 HOUR)
            AND connection_spec.data_direction = 1
            AND connection_spec.tls IS TRUE
            -- AND connection_spec.client_geolocation.continent_code = 'NA'
            AND web100_log_entry.snap.HCThruOctetsAcked >= 8192
            AND (web100_log_entry.snap.SndLimTimeRwin + web100_log_entry.snap.SndLimTimeCwnd + web100_log_entry.snap.SndLimTimeSnd) >= 9000000
            AND (web100_log_entry.snap.SndLimTimeRwin + web100_log_entry.snap.SndLimTimeCwnd + web100_log_entry.snap.SndLimTimeSnd) < 600000000
            AND connection_spec.client_geolocation.latitude IS NOT NULL
            AND connection_spec.client_geolocation.longitude IS NOT NULL
        GROUP BY
            metro_raw,
            lat,
            lon
    )
  )
)
WHERE
  row < lim
