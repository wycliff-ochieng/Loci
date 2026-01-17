-- name: Login :one
SELECT id,username,firstname,lastname,email,password_hash,created_at FROM Users WHERE email=$1 OR username=$2;