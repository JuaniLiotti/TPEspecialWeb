package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

func main() {

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

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
