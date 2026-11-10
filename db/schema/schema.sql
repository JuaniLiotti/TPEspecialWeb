CREATE TABLE producto(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    descripcion VARCHAR(255) NOT NULL,
    categoria VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL,
    precio INTEGER NOT NULL
);

CREATE TABLE cuero(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    tipo VARCHAR(255) NOT NULL,
    acabado VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL,
    espesor_mm DOUBLE PRECISION NOT NULL,
    descripcion VARCHAR(255) NOT NULL,
    precio INTEGER NOT NULL
);

CREATE TABLE usuario(
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- Imágenes: una entidad puede tener muchas imágenes (relación 1 a N).
-- El archivo se guarda en disco (carpeta uploads/) y acá solo su URL.
-- ON DELETE CASCADE: al borrar el producto/cuero se borran sus imágenes.
CREATE TABLE producto_imagen(
    id SERIAL PRIMARY KEY,
    producto_id INTEGER NOT NULL REFERENCES producto(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL
);

CREATE TABLE cuero_imagen(
    id SERIAL PRIMARY KEY,
    cuero_id INTEGER NOT NULL REFERENCES cuero(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL
);
