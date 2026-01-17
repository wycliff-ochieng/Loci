-- name: GetReplyForBroadcast :one
SELECT r.id, r.content, r.created_at, u.username, ST_Y(l.location::geometry)::float8 as lat, ST_X(l.location::geometry)::float8 as long
FROM replies r
JOIN users u ON r.user_id = u.id
JOIN loci l ON r.locus_id = l.id
WHERE r.id = $1;