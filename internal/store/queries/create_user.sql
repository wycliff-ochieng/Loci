-- name: CreateUser :exec
INSERT INTO Users(id,username,firstname,lastname,email,password_hash,created_at) VALUES($1,$2,$3,$4,$5,$6,$7);