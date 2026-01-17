-- name: GetLociInBounds :many
SELECT
    id,
    user_id,
    message,
    ST_Y(location::geometry)::float8 as lat,
    ST_X(location::geometry)::float8 as long,
    created_at,
    view_count,
    replies_count,
    visibility_score
FROM
    loci
WHERE
    ST_Within(location::geometry, ST_MakeEnvelope($1,$2,$3,$4, 4326))
    AND visibility_score >= 0.0
ORDER BY created_at DESC;

-- name: GetLocusLocation :one
SELECT id, ST_Y(location::geometry)::float8 as lat, ST_X(location::geometry)::float8 as long
FROM loci
WHERE id = $1;