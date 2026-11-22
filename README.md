# Italpiel S.A. - Interfaces dinámicas con HTMX

Italpiel S.A. es una plataforma para una tienda de la industria del cuero. La tienda se organiza en dos categorías:
* **Cueros:** cueros en sí mismos, como materia prima (vaquetas, cueros ovinos, etc.).
* **Productos terminados:** alfombras, bajo monturas, productos textiles, etc.

El proyecto es una aplicación web en Go que administra el catálogo. Permite gestionar tres entidades:
* **Cueros:** nombre, tipo, acabado, color, espesor (mm), descripción, precio y una o más imágenes.
* **Productos:** nombre, descripción, categoría, color, precio y una o más imágenes.
* **Usuarios:** nombre, email y contraseña (solo por la API).

Usa PostgreSQL como base de datos (con `sqlc` para el acceso tipado), `templ` para generar las vistas HTML en el servidor, HTMX para actualizarlas sin recargar la página, y Docker + Make para automatizar todo el entorno.

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

## TP6: interfaces dinámicas con HTMX

En el TP5 la página se armaba en el servidor con `templ`, pero cada alta o baja era un formulario clásico: `POST` + redirect `303` + recarga completa de la página (patrón Post/Redirect/Get). En el TP6 se incorpora [HTMX](https://htmx.org): los formularios y los botones envían el pedido por AJAX, el servidor responde **solo con el fragmento de HTML que cambió** (renderizado con los mismos componentes `templ`) y HTMX lo coloca en el lugar indicado. La página no se recarga y no se escribe JavaScript propio.

### 1. Integración de HTMX

En `views/layout.templ` se agregó el script de HTMX en el `<head>`, así queda disponible en todas las vistas que usan el `Layout`:

```html
<script src="https://unpkg.com/htmx.org@1.9.10"></script>
```

### 2. Creación sin recarga de página

**Lista** (`views/cuero_list.templ`, `views/producto_list.templ`): el contenedor de las tarjetas tiene un `id` (`lista-cueros` / `lista-productos`) para que HTMX pueda ubicarlo.

**Formulario** (`views/cuero_form.templ`, `views/producto_form.templ`): el `<form>` tiene los atributos de HTMX:

```html
<form id="form-cuero" enctype="multipart/form-data"
      hx-post="/cueros" hx-target="#lista-cueros" hx-swap="outerHTML">
```

**Handler** `POST /cueros` y `POST /productos` (`logic/web/cueros.go`, `logic/web/productos.go`): se eliminó el `http.Redirect`. Después de crear el registro (y guardar sus imágenes) se vuelve a consultar la lista completa y se responde con dos fragmentos:

1. La lista actualizada (`views.CueroList`), que va al `hx-target`.
2. El formulario vacío con `hx-swap-oob="true"` (`views.CueroFormOOB`). Es un *Out-of-Band swap*: HTMX lo saca de la respuesta, busca en la página el elemento con el mismo `id` (`form-cuero`) y lo reemplaza, así **se limpian los campos** después de un alta exitosa.

**Errores de validación** (precio inválido, archivo que no es imagen, etc.): el servidor responde **solo con el formulario**, con el mensaje de error y los datos que había cargado el usuario, y agrega dos headers de respuesta de HTMX:

| Header | Valor | Efecto |
|---|---|---|
| `HX-Retarget` | `#form-cuero` | El fragmento va al formulario en lugar de la lista |
| `HX-Reswap` | `outerHTML` | Reemplaza el formulario entero |

### 3. Borrado sin recarga de página

**Botón** (`views/cuero_list.templ`, `views/producto_list.templ`): el formulario `POST /cueros/{id}/eliminar` del TP5 se reemplazó por un botón con atributos de HTMX:

```html
<button hx-delete="/cueros/5"
        hx-target="closest article"
        hx-swap="outerHTML"
        hx-confirm="¿Estás seguro?">Eliminar</button>
```

### Rutas de la página

| Ruta | Descripción |
|---|---|
| `GET /` | Página principal completa (`?ver=cueros` o `?ver=productos`) |
| `POST /cueros`, `POST /productos` | Alta por HTMX (`multipart/form-data`). Responde con la lista y el formulario vacío (OOB), o con el formulario con el error |
| `DELETE /cueros/{id}`, `DELETE /productos/{id}` | Baja por HTMX. Responde vacío con `200` |

Las rutas `POST /cueros/{id}/eliminar` y `POST /productos/{id}/eliminar` del TP5 se eliminaron. La API JSON bajo `/api/` no cambió.

