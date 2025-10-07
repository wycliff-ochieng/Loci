-- name: CreateLoci :many
INSERT INTO loci(
    user_id,
    message,
    location )
VALUES
    ($1,$2,$3) 
RETURNING *;