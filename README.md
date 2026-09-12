# TPEspecialWeb

# Trabajo Práctico 1:

Italpiel S.A.

Nuestra página es una tienda de productos de la industria del cuero, en el que se mostrarán diversos productos como alfombras, bajo monturas, productos textiles, etc.
Cada producto debe tener nombre, material, categoría, imagen, color, tamaño y precio.
Inicialmente pensamos en hacer una página informativa sobre los productos para luego crear la tienda con distintas categorías.
Para ejecutarlo hay que correr el siguiente comando desde la carpeta del trabajo: go run main.go

El sitio web estara en localhost:8080/

# Trabajo Práctico 2:

## Ejecución:
Para probar el proyecto y correr el test completo, ejecute los siguientes comandos:
(Bash) 
make test

Este comando se encargará de:

Ejecutar todos los comandos previos al test:
    - Genera el código sqlc
    - Compila el programa
    - Crea la imagen de la BBDD utilizando el Dockerfile
    - Borra contenedores y volúmenes viejos 
    - Levanta la BBDD
    

Ejecutar el test:
    - De manera automatizada ejecuta los tests que se encuentran en db_test.go (Crear, obtener uno, obtener todos, actualizar, eliminar)

## Documentación:
En db/schema se encuentran los esquemas de creación de las tablas:
Este esquema representa una simplificación del dominio de la página, ya que contiene a los productos que serán exhibidos en la página, permitiendo guardar un ID, nombre de producto, descripción, categoría, color y precio. Y también a los usarios que tendrán la oportunidad de logearse, guardando ennla BBDD la id de usuario, el nombre, el mail y la clave de acceso. 

En db/queries se encuentran las sentencias SQL para realizar los métodos CRUD (cada nombre .sql representa la tabla correspondiente).
Esta carpeta es de utilidad, ya que a partir de ella se autogeneran las operaciones mediante el sqlc.

En db/generated se encuentrae el código GO autogenerado
Se comprueba que este código es correcto y funcional, ya que todas las operacions CRUD pasan el testing.

### Contenedor de BBDD:
Imagen: Custom tp-postgres, construida a partir del Dockerfile propio de la raíz del proyecto.
Nombre del contenedor: postgres-tp2
Variables del entorno definidas en Dockerfile: 
    ENV POSTGRES_USER=postgres
    ENV POSTGRES_PASSWORD=secreta
    ENV POSTGRES_DB=italpiel_db
Puerto: 5432 del host al 5432 del contenedor
Para inicializar el schema: Se define en el Dockerfile el comando:
COPY ./db/schema/schema.sql /docker-entrypoint-initdb.d/schema.sql
El cual copia el schema guardado en db/schema/schema.sql a /docker-entrypoint-initdb.d/schema.sql,
a partir de este comando, Postgres ejecuta de manera automática cualquier script .sql de ese directorio la primera vez
que arranca data directory vacío, de esta forma las tablas usuario y producto se crean sin pasos manuales adicionales.

### Volúmenes y persistencia real:
Por el momento no generamos persistencia real, ya que al eliminar los datos al final de cada test, generamos que los 
datos creados solo existan durante la ejecución de este.