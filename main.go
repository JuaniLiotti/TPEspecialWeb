package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"tpespecialweb/db"
	sqlc "tpespecialweb/db/sqlc"
	"tpespecialweb/logic/cueroCRUD"
	"tpespecialweb/logic/imagenes"
	"tpespecialweb/logic/productCRUD"
	"tpespecialweb/logic/userCRUD"
	"tpespecialweb/logic/web"

	_ "github.com/lib/pq"
)

func pathExists(urlPath string, staticDir string) bool {
	path := filepath.Join(staticDir, filepath.Clean(urlPath))
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	// si la url es un directorio, tiene que chequear que haya un index
	if info.IsDir() {
		indexPath := filepath.Join(path, "index.html")
		if _, err := os.Stat(indexPath); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Imprime en consola: Método HTTP, URL y IP remota
		log.Printf("[%s] %s - Desde IP: %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Pasa la petición al siguiente handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Conexión a la base de datos PostgreSQL
	conn, err := sql.Open("postgres", "postgres://postgres:secreta@localhost:5432/italpiel_db?sslmode=disable")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	conn.SetMaxOpenConns(10)
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Error al hacer ping a la base de datos:", err)
	}

	fmt.Println("Conexión a la BBDD exitosa")

	queries := sqlc.New(conn)

	productRepo := db.NewProductAdapter(queries)
	userRepo := db.NewUserAdapter(queries)
	cueroRepo := db.NewCueroAdapter(queries)

	prodHandler := productCRUD.NewProductHandler(productRepo)
	userHandler := userCRUD.NewUserHandler(userRepo)
	cueroHandler := cueroCRUD.NewCueroHandler(cueroRepo)

	// Páginas HTML renderizadas en el servidor con templ (TP5)
	webHandler := web.NewWebHandler(productRepo, cueroRepo)

	fmt.Println("BBDD Lista para operar")
	// Fin conexión BBDD

	port := ":8080"
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))

	// Página principal: la arma el servidor con templ (views.IndexPage).
	// "{$}" hace que matchee solo "/" exacto y no todo lo que empieza con "/".
	http.Handle("GET /{$}", loggingMiddleware(http.HandlerFunc(webHandler.Index)))

	// Formularios HTML clásicos (POST + redirect a GET /, patrón PRG)
	http.Handle("POST /cueros", loggingMiddleware(http.HandlerFunc(webHandler.CrearCuero)))
	http.Handle("POST /productos", loggingMiddleware(http.HandlerFunc(webHandler.CrearProducto)))
	http.Handle("POST /cueros/{id}/eliminar", loggingMiddleware(http.HandlerFunc(webHandler.EliminarCuero)))
	http.Handle("POST /productos/{id}/eliminar", loggingMiddleware(http.HandlerFunc(webHandler.EliminarProducto)))

	// Archivos estáticos (styles.css, about.html) y 404 para todo lo demás
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) // 405
			return
		}

		if pathExists(r.URL.Path, staticDir) {
			fileServer.ServeHTTP(w, r)
		} else {
			// Se escribe el archivo a mano: http.ServeFile redirige las rutas que
			// terminan en /index.html y dejaba el 404 con el cuerpo vacío
			pagina, err := os.ReadFile(filepath.Join(staticDir, "notfound.html"))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write(pagina)
		}
	}))

	// API REST JSON bajo /api/. StripPrefix quita "/api" para que cada handler
	// siga recibiendo rutas del tipo /products/{id}, /users/{id}, /cueros/{id}
	http.Handle("/api/products/", loggingMiddleware(http.StripPrefix("/api", prodHandler)))
	http.Handle("/api/users/", loggingMiddleware(http.StripPrefix("/api", userHandler)))
	http.Handle("/api/cueros/", loggingMiddleware(http.StripPrefix("/api", cueroHandler)))

	// Imágenes: patrones de Go 1.22+ con método y parámetros ({id}, {imgId}).
	// Son más específicos que "/api/products/", así que el mux los elige primero.
	http.Handle("POST /api/products/{id}/imagenes", loggingMiddleware(http.HandlerFunc(prodHandler.UploadImagenes)))
	http.Handle("DELETE /api/products/{id}/imagenes/{imgId}", loggingMiddleware(http.HandlerFunc(prodHandler.DeleteImagen)))
	http.Handle("POST /api/cueros/{id}/imagenes", loggingMiddleware(http.HandlerFunc(cueroHandler.UploadImagenes)))
	http.Handle("DELETE /api/cueros/{id}/imagenes/{imgId}", loggingMiddleware(http.HandlerFunc(cueroHandler.DeleteImagen)))

	// Archivos de imagen subidos, ej: /uploads/productos/3/a1b2c3.jpg
	http.Handle("GET /uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(imagenes.DirBase))))

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	error := http.ListenAndServe(port, nil)
	if error != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", error)
	}
}
