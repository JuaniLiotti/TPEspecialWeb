# Italpiel S.A. - Tienda de cueros

Italpiel S.A. es una plataforma para una tienda de la industria del cuero. La tienda se organiza en dos categorías:
* **Cueros:** cueros en sí mismos, como materia prima (vaquetas, cueros ovinos, etc.).
* **Productos terminados:** alfombras, bajo monturas, productos textiles, etc.

El proyecto es una aplicación web en Go que administra el catálogo. Permite gestionar tres entidades:
* **Cueros:** nombre, tipo, acabado, color, espesor (mm), descripción, precio y una o más imágenes.
* **Productos:** nombre, descripción, categoría, color, precio y una o más imágenes.
* **Usuarios:** nombre, email y contraseña (solo por la API).

Usa PostgreSQL como base de datos (con `sqlc` para el acceso tipado), `templ` para generar las vistas HTML en el servidor, y Docker + Make para automatizar todo el entorno.

## Cómo ejecutarlo

**Requisitos:** Go, [`sqlc`](https://docs.sqlc.dev/en/latest/overview/install.html), Docker, `make` y `curl` (y `jq` solo para `make test-api`).
No hace falta instalar `templ`: está declarado como *tool* en `go.mod` y se ejecuta con `go tool templ`.

```bash
git clone <url-del-repo>
cd TPEspecialWeb
make run-app
```

Después abrir **http://localhost:8080**. Para cortar la aplicación: **Ctrl+C**.

`make run-app` hace todo lo necesario, en este orden:

1. **Tareas previas**
   1. Genera el código de acceso a datos con `sqlc` (`db/sqlc/`).
   2. Genera los componentes de vista con `templ` (`views/*_templ.go`).
   3. Compila el servidor (`api-server`).
   4. Borra el contenedor de PostgreSQL de una corrida anterior, junto con sus volúmenes.
   5. Construye la imagen y levanta un contenedor nuevo, que carga `db/schema/schema.sql`.
   6. Espera a que la base esté aceptando conexiones (`pg_isready`).
2. **Ejecuta la aplicación** en primer plano, en el puerto 8080.
3. **Tareas posteriores** (al cortar con Ctrl+C, o si la aplicación se cae): borra el contenedor y sus volúmenes, y las imágenes subidas en `uploads/`. Cada corrida arranca con la base vacía, así que esas imágenes quedarían huérfanas.

> El código generado (`db/sqlc/` y `views/*_templ.go`) no se versiona: lo genera el Makefile en cada corrida.

### Todos los comandos

`make` (o `make help`) lista los comandos disponibles:

| Comando | Qué hace |
|---|---|
| `make run-app` | Tareas previas + ejecuta la app + tareas posteriores (`make run` es un alias) |
| `make test` | Tests de la base de datos (`go test`) con una BBDD nueva, y limpia al terminar |
| `make test-api` | Tests de la API JSON (`requests.sh`) con la app en segundo plano, y limpia al terminar |
| `make build-app` | Genera el código (sqlc + templ) y compila, sin levantar nada |
| `make generate` | Solo genera el código de sqlc y templ |
| `make db` | Levanta una BBDD nueva y vacía en Docker |
| `make clean` | Limpia todo: servidor en segundo plano, BBDD, imágenes subidas y binario |

## Uso de la página (TP5: renderizado en el servidor)

La página principal se genera en el servidor con `templ`. No usa JavaScript: el servidor arma el HTML completo con los datos de la base y los formularios se envían con un POST clásico.

- **Selector Cueros / Productos terminados:** son dos links (`/?ver=cueros` y `/?ver=productos`). El servidor renderiza solo la sección elegida y consulta en la base solo esa lista.
- **Alta:** cada sección tiene un formulario (con imágenes opcionales) que hace `POST /cueros` o `POST /productos`. Si los datos son válidos, el servidor guarda el registro y responde con un redirect `303` a la página principal: es el patrón **Post/Redirect/Get**, así que recargar la página no vuelve a enviar el formulario. Si hay un error (precio inválido, archivo que no es imagen, etc.), se vuelve a mostrar el formulario con el mensaje y con los datos cargados.
- **Eliminar:** cada tarjeta tiene un botón que hace `POST /cueros/{id}/eliminar` o `POST /productos/{id}/eliminar`, y también redirige a la página principal.

| Ruta | Descripción |
|---|---|
| `GET /` | Página principal (`?ver=cueros` o `?ver=productos`) |
| `POST /cueros`, `POST /productos` | Alta desde el formulario (`multipart/form-data`) |
| `POST /cueros/{id}/eliminar`, `POST /productos/{id}/eliminar` | Baja desde la tarjeta |
| `GET /about.html` | Página estática "Sobre nosotros" |
| `GET /uploads/...` | Archivos de imagen subidos |

Las reglas de validación son las mismas para los formularios y para la API (funciones `Validate` de `logic/cueroCRUD` y `logic/productCRUD`).

### Estructura del proyecto

```
main.go                 Conexión a la BBDD y registro de rutas
views/                  Componentes templ (layout, listas, formularios, página principal)
logic/web/              Handlers de las páginas HTML y de los formularios
logic/cueroCRUD/        Lógica y API JSON de cueros
logic/productCRUD/      Lógica y API JSON de productos
logic/userCRUD/         Lógica y API JSON de usuarios
logic/imagenes/         Subida, validación y borrado de imágenes
db/                     Adaptadores entre la lógica y el código de sqlc
db/schema/, db/queries/ Esquema y consultas SQL (entrada de sqlc)
static/                 styles.css, about.html y notfound.html
```

## Tests automatizados

### Tests de la base de datos: `make test`
1. Genera el código con `sqlc` y `templ`.
2. Levanta un contenedor de PostgreSQL con el esquema inicial.
3. Ejecuta `go test`, que recorre el ciclo CRUD completo de usuarios, productos, cueros y sus imágenes (incluyendo el borrado en cascada).
4. Elimina el contenedor y sus volúmenes.

### Tests de la API: `make test-api`
1. Genera el código y compila la aplicación (`api-server`).
2. Levanta un contenedor de PostgreSQL con el esquema inicial.
3. Arranca el servidor en segundo plano (puerto 8080) y espera a que responda.
4. Ejecuta `requests.sh`, que hace peticiones HTTP reales (POST, GET, PUT, DELETE) contra la API, incluyendo la subida de imágenes de prueba (carpeta `testdata/`) y casos de error (`400`, `404`).
5. Mata el servidor y elimina el contenedor, sus volúmenes y las imágenes subidas.

## API JSON

Además de las páginas HTML, el servidor mantiene la API REST de los TPs anteriores. Todos sus endpoints están bajo el prefijo `/api/`.

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

La subida se hace como `multipart/form-data`, con los archivos en el campo `imagenes`:
```bash
curl -X POST http://localhost:8080/api/products/1/imagenes \
    -F "imagenes=@foto1.jpg" -F "imagenes=@foto2.png"
```

Reglas de validación (también se aplican a las imágenes de los formularios HTML):
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
