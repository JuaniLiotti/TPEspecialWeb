# TPEspecialWeb

# Trabajo Práctico 1:

Italpiel S.A.

Nuestra página es una tienda de productos de la industria del cuero, en el que se mostrarán diversos productos como alfombras, bajo monturas, productos textiles, etc.
Cada producto debe tener nombre, material, categoría, imagen, color, tamaño y precio.
Inicialmente pensamos en hacer una página informativa sobre los productos para luego crear la tienda con distintas categorías.
Para ejecutarlo hay que correr el siguiente comando desde la carpeta del trabajo: `go run main.go`

El sitio web estará en localhost:8080/

---

# Trabajo Práctico 2:

## Ejecución:
Para probar el proyecto y correr el test completo de la base de datos, ejecute el siguiente comando:
(Bash) 
`make test`

Este comando se encargará de:

Ejecutar todos los comandos previos al test:
    - Genera el código sqlc
    - Crea la imagen de la BBDD utilizando el Dockerfile
    - Borra contenedores y volúmenes viejos 
    - Levanta la BBDD
    
Ejecutar el test:
    - De manera automatizada ejecuta los tests que se encuentran en `db_test.go` (Crear, obtener uno, obtener todos, actualizar, eliminar)

Limpiar BBDD post test: 
    - A partir del comando: `docker rm -f -v $(CONTAINER_NAME)` se eliminan los datos creados para la prueba, dejando la base de datos completamente vacía.

## Documentación:
En `db/schema` se encuentran los esquemas de creación de las tablas:
Este esquema representa una simplificación del dominio de la página, ya que contiene a los productos que serán exhibidos en la página, permitiendo guardar un ID, nombre de producto, descripción, categoría, color y precio. Y también a los usuarios que tendrán la oportunidad de logearse, guardando en la BBDD la id de usuario, el nombre, el mail y la clave de acceso. 

En `db/queries` se encuentran las sentencias SQL para realizar los métodos CRUD (cada nombre .sql representa la tabla correspondiente).
Esta carpeta es de utilidad, ya que a partir de ella se autogeneran las operaciones mediante el sqlc.

En `db/sqlc/` se encuentra el código GO autogenerado (esta carpeta está ignorada en el repositorio ya que se autogenera en la compilación).
Se comprueba que este código es correcto y funcional, ya que todas las operaciones CRUD pasan el testing.

### Contenedor de BBDD:
Imagen: Custom tp-postgres, construida a partir del Dockerfile propio de la raíz del proyecto.
Nombre del contenedor: postgres-tp2
Variables del entorno definidas en Dockerfile: 
    ENV POSTGRES_USER=postgres
    ENV POSTGRES_PASSWORD=secreta
    ENV POSTGRES_DB=italpiel_db
Puerto: 5432 del host al 5432 del contenedor
Para inicializar el schema: Se define en el Dockerfile el comando:
`COPY ./db/schema/schema.sql /docker-entrypoint-initdb.d/schema.sql`
El cual copia el schema guardado en `db/schema/schema.sql` a `/docker-entrypoint-initdb.d/schema.sql`, a partir de este comando, Postgres ejecuta de manera automática cualquier script .sql de ese directorio la primera vez que arranca data directory vacío, de esta forma las tablas usuario y producto se crean sin pasos manuales adicionales.

### Volúmenes y persistencia real:
Por el momento no generamos persistencia real, ya que al eliminar los datos al final de cada test, generamos que los datos creados solo existan durante la ejecución de este.

---

# Trabajo Práctico 3:

En esta etapa conectamos nuestra capa de datos (generada con sqlc) con una API RESTful desarrollada en Go, exponiendo las operaciones CRUD tanto para Productos como para Usuarios mediante endpoints HTTP.

## Ejecución de la API y Pruebas Automatizadas

Para automatizar todo el flujo de pruebas, hemos incorporado un nuevo comando en nuestro `Makefile`:

`make test-api`

Este comando ejecuta las siguientes acciones secuencialmente:
1. **Tareas Previas:** Genera el código de `sqlc`, compila la aplicación Go (`api-server`), levanta el contenedor de PostgreSQL y luego arranca el servidor web en segundo plano (puerto 8080).
2. **Pruebas de API:** Ejecuta automáticamente el script `requests.sh`, el cual utiliza `curl` para realizar peticiones HTTP reales (POST, GET, PUT, DELETE) contra los endpoints de Productos y Usuarios.
3. **Limpieza:** Al finalizar, mata el proceso del servidor web (`api-server`) y elimina los contenedores y volúmenes de Docker, dejando el entorno limpio.

## Estructura y Arquitectura del Código

El código ha sido refactorizado:
* **`logic/productCRUD/` y `logic/userCRUD/`**: Se crearon paquetes separados para manejar la lógica de negocio y los controladores HTTP (handlers).
* **Inyección de Dependencias:** Los controladores ya no utilizan variables globales. Reciben la conexión a la base de datos a través de una estructura (Struct) en su constructor, lo que permite mayor escalabilidad y facilidad para el testing.
* **Middleware:** Se implementó un middleware de logging en `main.go` que registra el método HTTP, la URL y la IP de origen de cada petición en la consola.

## Endpoints Disponibles

La API valida la entrada de datos (ej. devolviendo `400 Bad Request` si faltan campos obligatorios) y maneja códigos de estado HTTP estándar (`201 Created`, `204 No Content`, `404 Not Found`).

### Productos (`/products/`)
- `GET /products/`: Lista todos los productos en formato JSON.
- `GET /products/{id}`: Obtiene un producto por su ID.
- `POST /products/`: Crea un producto. Validaciones: El nombre no puede estar vacío y el precio debe ser > 0.
- `PUT /products/{id}`: Actualiza todos los campos de un producto por su ID.
- `DELETE /products/{id}`: Elimina un producto por su ID.

### Usuarios (`/users/`)
- `GET /users/`: Lista todos los usuarios en formato JSON.
- `GET /users/{id}`: Obtiene un usuario por su ID.
- `POST /users/`: Crea un usuario. Validaciones: Nombre, email y contraseña son obligatorios.
- `PUT /users/{id}`: Actualiza los datos de un usuario por su ID.
- `DELETE /users/{id}`: Elimina un usuario por su ID.
