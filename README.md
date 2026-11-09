# Italpiel S.A. - API RESTful

Italpiel S.A. es una plataforma para una tienda de la industria del cuero. La tienda se organiza en dos categorías:
* **Cueros:** cueros en sí mismos, como materia prima (vaquetas, cueros ovinos, etc.).
* **Productos terminados:** alfombras, bajo monturas, productos textiles, etc.

Este proyecto implementa la API RESTful backend en Go que administra el sistema. Permite gestionar tres entidades:
* **Cueros:** Almacenando nombre, tipo, acabado, color, espesor (mm), descripción, precio y una o más imágenes.
* **Productos:** Almacenando nombre, descripción, categoría, color, precio y una o más imágenes.
* **Usuarios:** Gestionando el acceso mediante nombre, email y contraseña.

El sistema utiliza PostgreSQL como motor de base de datos (gestionado mediante `sqlc` para type-safety) y despliega su infraestructura de datos y pruebas utilizando Docker y automatización vía Make.

## Ejecución de las Pruebas Automatizadas

El proyecto incluye un `Makefile` diseñado para levantar el entorno completo y ejecutar las pruebas automáticamente.

**Requisitos previos:** Go, `sqlc`, Docker, `make`, `curl` y `jq`.

### Tests de la base de datos
`make test`

1. Genera el código de acceso a datos con `sqlc`.
2. Levanta un contenedor de PostgreSQL con el esquema inicial.
3. Ejecuta `go test`, que recorre el ciclo CRUD completo de usuarios, productos, cueros y sus imágenes (incluyendo el borrado en cascada).
4. Elimina el contenedor, limpiando el entorno.

### Tests de la API
`make test-api`

Este comando ejecuta secuencialmente:
1. Genera el código de acceso a datos con `sqlc`.
2. Levanta un contenedor de PostgreSQL con el esquema inicial.
3. Compila la aplicación Go (`api-server`).
4. Arranca el servidor web en segundo plano (puerto 8080).
5. Ejecuta automáticamente el script `requests.sh`, realizando peticiones HTTP reales (POST, GET, PUT, DELETE) contra la API, incluyendo la subida de imágenes de prueba (carpeta `testdata/`) y casos de error (`400`, `404`).
6. Mata el proceso del servidor y elimina los contenedores, limpiando el entorno.

## Endpoints Disponibles

Todos los endpoints de la API están bajo el prefijo `/api/`.

La API valida la entrada de datos mediante DTOs propios y funciones de validación en la capa de lógica, y maneja códigos de estado HTTP estándar (`200 OK`, `201 Created`, `204 No Content`, `400 Bad Request`, `404 Not Found`).

### Cueros (`/api/cueros/`)
- `GET /api/cueros/`: Lista todos los cueros.
- `GET /api/cueros/{id}`: Obtiene un cuero por ID.
- `POST /api/cueros/`: Crea un cuero.
- `PUT /api/cueros/{id}`: Actualiza un cuero por ID.
- `DELETE /api/cueros/{id}`: Elimina un cuero.

Ejemplo de cuerpo JSON:
```json
{"nombre": "Vaqueta Natural", "tipo": "Vacuno", "acabado": "Curtido vegetal", "color": "Suela", "espesor_mm": 2.5, "descripcion": "Cuero entero", "precio": 120000}
```

### Productos (`/api/products/`)
- `GET /api/products/`: Lista todos los productos.
- `GET /api/products/{id}`: Obtiene un producto por ID.
- `POST /api/products/`: Crea un producto.
- `PUT /api/products/{id}`: Actualiza un producto por ID.
- `DELETE /api/products/{id}`: Elimina un producto.

Ejemplo de cuerpo JSON:
```json
{"nombre": "Cartera de Cuero", "descripcion": "100% cuero natural", "categoria": "Accesorios", "color": "Marrón", "precio": 45000}
```

### Imágenes de cueros y productos
Cada cuero y cada producto puede tener una o más imágenes. Los archivos se guardan en disco, en la carpeta `uploads/`, y en la base de datos solo se guarda su URL. Las respuestas `GET` de cueros y productos incluyen el campo `imagenes`; la primera imagen de la lista es la portada:
```json
"imagenes": [{"id": 1, "url": "/uploads/productos/1/b88b6d9c91b95cd82fca167c343f4333.png"}]
```

- `POST /api/products/{id}/imagenes`: Sube una o más imágenes al producto. Responde `201` con las imágenes creadas.
- `DELETE /api/products/{id}/imagenes/{imgId}`: Elimina una imagen del producto (fila y archivo). Responde `204`.
- `POST /api/cueros/{id}/imagenes` y `DELETE /api/cueros/{id}/imagenes/{imgId}`: Lo mismo para cueros.
- `GET /uploads/...`: Sirve los archivos de imagen.

La subida se hace como `multipart/form-data`, con los archivos en el campo `imagenes`:
```bash
curl -X POST http://localhost:8080/api/products/1/imagenes \
    -F "imagenes=@foto1.jpg" -F "imagenes=@foto2.png"
```

Reglas de validación:
- Solo se aceptan imágenes JPG, PNG o WebP. El tipo se detecta por el contenido del archivo, no por su extensión.
- Máximo 10 MB por request.
- El servidor genera el nombre de cada archivo (nunca usa el nombre original).

Al eliminar un cuero o un producto se borran también sus imágenes: las filas, por `ON DELETE CASCADE`, y los archivos de su carpeta.

### Usuarios (`/api/users/`)
- `GET /api/users/`: Lista todos los usuarios.
- `GET /api/users/{id}`: Obtiene un usuario por ID.
- `POST /api/users/`: Crea un usuario.
- `PUT /api/users/{id}`: Actualiza un usuario por ID.
- `DELETE /api/users/{id}`: Elimina un usuario.
