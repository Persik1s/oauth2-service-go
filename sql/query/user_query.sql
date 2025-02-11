-- name: CreateUser :one
INSERT INTO "users"(username,password,email,createAt) VALUES($1,$2,$3,$4) RETURNING id;

-- name: GetUser :one
SELECT * FROM "users" WHERE username = $1;

-- name: CreateUserRole :exec
INSERT INTO "user_roles"(user_id,role_id) VALUES($1,$2);

-- name: GetUserRole :one
SELECT * FROM  "user_roles" WHERE user_id = $1;