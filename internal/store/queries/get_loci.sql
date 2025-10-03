-- name: GetLociInBounds :many
SELECT
    id,
    user_id,
    message
    location,
    created_at,
    view_count,
    replies_count,
    visibility_score
FROM
    loci
WHERE
    ST_Within(location::geometry, ST_MakeEnvelope($1,$2,$3,$4,$5, 4326))
    AND visibility_score > 0.1
ORDER BY created_at DESC;