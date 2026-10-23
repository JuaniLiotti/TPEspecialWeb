# Italpiel S.A. - API RESTful

Italpiel S.A. es una plataforma para una tienda de productos de la industria del cuero (alfombras, bajo monturas, productos textiles, etc.). 

Este proyecto implementa la API RESTful backend en Go que administra el sistema. Permite gestionar dos entidades principales:
* **Productos:** Almacenando nombre, descripción, categoría, color y precio.
* **Usuarios:** Gestionando el acceso mediante nombre, email y contraseña.

El sistema utiliza PostgreSQL como motor de base de datos (gestionado mediante `sqlc` para type-safety) y despliega su infraestructura de datos y pruebas utilizando Docker y automatización vía Make.

## Ejecución de la API y Pruebas Automatizadas

El proyecto incluye un `Makefile` diseñado para levantar el entorno completo y ejecutar las pruebas de los endpoints HTTP (CRUD) automáticamente.

**Requisitos previos:** Docker, `make`, `curl` y `jq`.

Para probar el proyecto, ejecute en la raíz:
`make test-api`

Este comando ejecuta secuencialmente:
1. Genera el código de acceso a datos con `sqlc`.
2. Compila la aplicación Go (`api-server`).
3. Levanta un contenedor de PostgreSQL con el esquema inicial.
4. Arranca el servidor web en segundo plano (puerto 8080).
5. Ejecuta automáticamente el script `requests.sh`, realizando peticiones HTTP reales (POST, GET, PUT, DELETE) contra la API.
6. Mata el proceso del servidor y elimina los contenedores, limpiando el entorno.

## Endpoints Disponibles

La API valida la entrada de datos mediante DTOs propios en la capa de lógica y maneja códigos de estado HTTP estándar (`200 OK`, `201 Created`, `204 No Content`, `400 Bad Request`, `404 Not Found`).

### Productos (`/products/`)
- `GET /products/`: Lista todos los productos.
- `GET /products/{id}`: Obtiene un producto por ID.
- `POST /products/`: Crea un producto.
- `PUT /products/{id}`: Actualiza un producto por ID.
- `DELETE /products/{id}`: Elimina un producto.

### Usuarios (`/users/`)
- `GET /users/`: Lista todos los usuarios.
- `GET /users/{id}`: Obtiene un usuario por ID.
- `POST /users/`: Crea un usuario.
- `PUT /users/{id}`: Actualiza un usuario por ID.
- `DELETE /users/{id}`: Elimina un usuario.