-- name: CreateProducto :one
INSERT INTO producto (nombre, descripcion, categoria, color, precio)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, nombre, descripcion, categoria, color, precio;

-- name: GetProducto :one
SELECT nombre, descripcion, categoria, color, precio
FROM producto
WHERE id = $1;

-- name: GetAllProductos :many
SELECT nombre, descripcion, categoria, color, precio
FROM producto;

-- name: UpdateProducto :one
UPDATE producto
SET nombre = $2, descripcion = $3, categoria = $4, color = $5, precio = $6
WHERE id = $1
RETURNING nombre, descripcion, categoria, color, precio;

-- name: DeleteProducto :exec
DELETE FROM producto
WHERE id = $1;