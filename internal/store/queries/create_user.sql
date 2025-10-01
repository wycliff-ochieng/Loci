-- CreateUser :one
INSERT INTO Users(id,username,firstname,lastname,email,hashedpassword,createat,updatedat) VALUES($1,$2,$3,$4,$5,$6,$7)