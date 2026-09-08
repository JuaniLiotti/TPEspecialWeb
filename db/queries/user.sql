-- name: CreateUsuario :one
INSERT INTO usuario (nombre, email, password)
VALUES ($1, $2, $3)
RETURNING id, nombre, email;

-- name: GetUsuario :one
SELECT id, nombre, email
FROM usuario
WHERE id = $1;

-- name: GetAllUsuarios :many
SELECT id, nombre, email
FROM usuario;

-- name: UpdateUsuario :one
UPDATE usuario
SET nombre = $2, email = $3, password = $4
WHERE id = $1
RETURNING nombre, email, password;

-- name: DeleteUsuario :exec
DELETE FROM usuario
WHERE id = $1;
