-- Imágenes de productos

-- name: CreateProductoImagen :one
INSERT INTO producto_imagen (producto_id, url)
VALUES ($1, $2)
RETURNING id, producto_id, url;

-- name: GetImagenesByProducto :many
SELECT id, producto_id, url
FROM producto_imagen
WHERE producto_id = $1
ORDER BY id;

-- name: GetAllProductoImagenes :many
SELECT id, producto_id, url
FROM producto_imagen
ORDER BY producto_id, id;

-- name: GetProductoImagen :one
SELECT id, producto_id, url
FROM producto_imagen
WHERE id = $1 AND producto_id = $2;

-- name: DeleteProductoImagen :exec
DELETE FROM producto_imagen
WHERE id = $1 AND producto_id = $2;


-- Imágenes de cueros

-- name: CreateCueroImagen :one
INSERT INTO cuero_imagen (cuero_id, url)
VALUES ($1, $2)
RETURNING id, cuero_id, url;

-- name: GetImagenesByCuero :many
SELECT id, cuero_id, url
FROM cuero_imagen
WHERE cuero_id = $1
ORDER BY id;

-- name: GetAllCueroImagenes :many
SELECT id, cuero_id, url
FROM cuero_imagen
ORDER BY cuero_id, id;

-- name: GetCueroImagen :one
SELECT id, cuero_id, url
FROM cuero_imagen
WHERE id = $1 AND cuero_id = $2;

-- name: DeleteCueroImagen :exec
DELETE FROM cuero_imagen
WHERE id = $1 AND cuero_id = $2;
