-- 1. Rekap Jam Pemakaian Ruangan pada Bulan Berjalan
SELECT
    r.id AS room_id,
    r.name,
    COALESCE(
        ROUND(
            SUM(EXTRACT(EPOCH FROM (b.end_time - b.start_time)) / 3600.0)::numeric,
            2
        ),
        0
    ) AS total_duration_hours
FROM rooms r
LEFT JOIN bookings b
    ON b.room_id = r.id
    AND b.status = 'CONFIRMED'
    AND DATE_TRUNC('month', b.start_time AT TIME ZONE 'UTC') = DATE_TRUNC('month', NOW() AT TIME ZONE 'UTC')
GROUP BY r.id, r.name
ORDER BY r.id ASC;

-- 2. Pencarian Ruangan Tersedia pada Rentang Waktu Tertentu
SELECT
    r.id,
    r.name,
    r.capacity,
    r.is_active
FROM rooms r
WHERE r.is_active = TRUE
  AND NOT EXISTS (
      SELECT 1
      FROM bookings b
      WHERE b.room_id = r.id
        AND b.status = 'CONFIRMED'
        AND b.start_time < '2026-12-01T12:00:00Z'::timestamptz
        AND b.end_time > '2026-12-01T10:00:00Z'::timestamptz
  )
ORDER BY r.id ASC;
