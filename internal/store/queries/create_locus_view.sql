-- name: CreateView :one
INSERT INTO locus_views (
    user_id,
    locus_id,
    viewed_at
) VALUES (
    $1, $2, NOW() -- Use NOW() here so you don't have to pass it from Go
)
ON CONFLICT (user_id, locus_id) DO NOTHING
RETURNING user_id, locus_id, viewed_at;