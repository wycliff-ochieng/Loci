-- name: IncrementViewCount :exec
UPDATE loci 
SET view_count = view_count + 1 
WHERE id = $1;