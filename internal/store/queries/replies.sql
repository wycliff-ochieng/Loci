-- name: CreateReply :one
INSERT INTO replies (locus_id, user_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: IncrementReplyCount :exec
UPDATE loci 
SET replies_count = replies_count + 1
WHERE id = $1;

-- name: GetRepliesByLocus :many
SELECT 
	r.id,
	r.locus_id,
	r.user_id,
	r.content,
	r.created_at
FROM replies r
WHERE r.locus_id = $1
ORDER BY r.created_at DESC;