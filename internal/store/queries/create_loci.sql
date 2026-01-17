-- name: CreateLoci :many
INSERT INTO loci(
    id,
    user_id,
    message,
    location )
VALUES
    ($1,$2,$3,
    ST_MakePoint($4, $5)::geography) 
RETURNING
    id, 
    user_id, 
    message, 
    ST_Y(location::geometry)::float8 as lat,
    ST_X(location::geometry)::float8 as long,
    created_at, 
    view_count, 
    replies_count, 
    visibility_score;