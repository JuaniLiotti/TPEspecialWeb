package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	db "tpespecialweb/db/sqlc"

	_ "github.com/lib/pq"
)

func pathExists(urlPath string, staticDir string) bool {
	path := filepath.Join(staticDir, filepath.Clean(urlPath))
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	// si la url es un directorio, tiene que checkear que haya un index
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

func GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var listaProductos []db.GetAllProductosRow

	mu.RLock()
	listaProductos, err = db.GetAllProductos(ctx) // arreglar esto, tiene que pedir la lista de productos
	mu.RUnlock()

	// mandejar error

	json.NewEncoder(w).Encode(listaProductos)
}

func GetProductById(w http.ResponseWriter, r *http.Request, idStr string) {

	w.Header().Set("Content-Type", "application/json")

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	var productoEncontrado db.GetProductoRow

	mu.RLock()
	productoEncontrado, err = db.GetProducto(id) // tiene que pedir el producto dependiendo del id
	mu.RUnlock()

	// manejar error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(productoEncontrado)
}

func CreateProducto(w http.ResponseWriter, r *http.Request) {
	var nuevoProducto db.CreateProductoParams

	err1 := json.NewDecoder(r.Body).Decode(&nuevoProducto)

	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	// validar entrada

	mu.Lock()
	err := db.CreateProducto(nuevoProducto) // hacer que ande
	mu.Unlock()

	// checkeo error

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoProducto)
}

func DeleteProductoById(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	err2 := db.DeleteProducto(, id) // ya sabes

	// manejar error

	// Respuesta 204 sin cuerpo
	w.WriteHeader(http.StatusNoContent)
}

func UpdateProductoById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var nuevoProducto db.UpdateProductoParams
	err1 := json.NewDecoder(r.Body).Decode(&nuevoProducto)
	
	// validar
	
	mu.Lock()
	defer mu.Unlock()

	err2 := db.UpdateProducto(, nuevoProductom) // que ande

	// manejar error

	pActualizado.ID = id
	slice[indexEncontrado] = pActualizado

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pActualizado)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")

	switch r.Method {
	case http.MethodGet:
		if idStr == "" {
			GetProducts(w, r)
		} else {
			GetProductById(w, r, idStr)
		}
	case http.MethodPost:
		CreateProducto(w, r)
	case http.MethodPut:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para actualizar", http.StatusBadRequest)
			return
		}
		UpdateProductoById(w, r)
	case http.MethodDelete:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para eliminar", http.StatusBadRequest)
			return
		}
		DeleteProductoById(w, r, idStr)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
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

	db.New(conn)

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

	//handler CRUD de productos
	handlerBase := http.HandlerFunc(productsHandler)

	http.Handle("/products/", loggingMiddleware(handlerBase))

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	error := http.ListenAndServe(port, nil)
	if error != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", error)
	}
}
