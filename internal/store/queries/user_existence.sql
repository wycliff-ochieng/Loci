-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1);
