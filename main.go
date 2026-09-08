package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tpespecialweb/db/generated"

	_ "github.com/lib/pq"
)

func main() {
	// Conexión a la base de datos PostgreSQL
	db, err := sql.Open("postgres", "postgres://postgres:secreta@localhost:5432/italpiel_db?sslmode=disable")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Error al hacer ping a la base de datos:", err)
	}

	fmt.Println("Conexión a la BBDD exitosa")

	queries := generated.New(db)

	fmt.Println("BBDD Lista para operar")
	// Fin conexión BBDD

	port := ":8080"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	error := http.ListenAndServe(port, nil)
	if error != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", error)
	}
}
