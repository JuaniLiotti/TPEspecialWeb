-- name: CreateCuero :one
INSERT INTO cuero (nombre, tipo, acabado, color, espesor_mm, descripcion, precio)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, nombre, tipo, acabado, color, espesor_mm, descripcion, precio;

-- name: GetCuero :one
SELECT id, nombre, tipo, acabado, color, espesor_mm, descripcion, precio
FROM cuero
WHERE id = $1;

-- name: GetAllCueros :many
SELECT id, nombre, tipo, acabado, color, espesor_mm, descripcion, precio
FROM cuero
ORDER BY nombre;

-- name: UpdateCuero :one
UPDATE cuero
SET nombre = $2, tipo = $3, acabado = $4, color = $5, espesor_mm = $6, descripcion = $7, precio = $8
WHERE id = $1
RETURNING id, nombre, tipo, acabado, color, espesor_mm, descripcion, precio;

-- name: DeleteCuero :exec
DELETE FROM cuero
WHERE id = $1;
