package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	db "tpespecialweb/db/sqlc"
	"tpespecialweb/logic/productCRUD"
	"tpespecialweb/logic/userCRUD"

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

	queries := db.New(conn)

	prodHandler := productCRUD.NewProductHandler(queries)
	userHandler := userCRUD.NewUserHandler(queries)

	fmt.Println("BBDD Lista para operar")
	// Fin conexión BBDD

	port := ":8080"
	staticDir := "./static"
	fileServer := http.FileServer(http.Dir(staticDir))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed) // 405
			return
		}

		if pathExists(r.URL.Path, staticDir) {
			fileServer.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, filepath.Join(staticDir, "notfound.html"))
		}
	}))

	// handler CRUD de productos
	http.Handle("/products/", loggingMiddleware(http.HandlerFunc(prodHandler.ServeHTTP)))
	http.Handle("/users/", loggingMiddleware(http.HandlerFunc(userHandler.ServeHTTP)))

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	error := http.ListenAndServe(port, nil)
	if error != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", error)
	}
}
