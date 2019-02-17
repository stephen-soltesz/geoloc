SELECT
  lat, lon,
  CASE
      WHEN metro = "sjc" THEN "nuq"
      ELSE metro
  END as metro
FROM (
      SELECT
        REGEXP_EXTRACT(connection_spec.server_hostname, r"mlab\d.([a-z]{3})..") as metro,
        connection_spec.client_geolocation.latitude as lat,
        connection_spec.client_geolocation.longitude as lon 
      FROM
        `measurement-lab.base_tables.ndt`
      WHERE
        _PARTITIONTIME >= TIMESTAMP("2019-02-14")
        AND connection_spec.data_direction = 1
        -- AND connection_spec.tls is True
        AND web100_log_entry.snap.HCThruOctetsAcked >= 8192
        AND (web100_log_entry.snap.SndLimTimeRwin +
          web100_log_entry.snap.SndLimTimeCwnd +
          web100_log_entry.snap.SndLimTimeSnd) >= 9000000
        AND (web100_log_entry.snap.SndLimTimeRwin +
          web100_log_entry.snap.SndLimTimeCwnd +
          web100_log_entry.snap.SndLimTimeSnd) < 600000000
        --AND web100_log_entry.snap.CongSignals > 0
        --AND (web100_log_entry.snap.State = 1 OR
        --  (web100_log_entry.snap.State >= 5 AND
        --  web100_log_entry.snap.State <= 11))
        AND connection_spec.client_geolocation.latitude is not NULL
        AND connection_spec.client_geolocation.longitude is not NULL
      GROUP BY
        metro, lat, lon
)
-- WHERE
--   row < 10000