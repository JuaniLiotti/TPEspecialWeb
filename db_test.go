package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	db "tpespecialweb/db/sqlc"
	sqlc "tpespecialweb/db/sqlc"

	_ "github.com/lib/pq"
)

func TestQueries_CRUD(t *testing.T) {
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
	ctx := context.Background()

	fmt.Println("BBDD Lista para operar")

	// tests para usuario
	var userID int32
	nombreUsuarioTest := "usuario test"
	emailTest := "test@test.com"
	passwordTest := "contra test"

	t.Run("CreateUsuario", func(t *testing.T) {
		user, err := queries.CreateUsuario(ctx, sqlc.CreateUsuarioParams{
			Nombre:   nombreUsuarioTest,
			Email:    emailTest,
			Password: passwordTest,
		})
		if err != nil {
			t.Fatalf("Error al crear usuario %v", err)
		}
		userID = user.ID
	})

	t.Run("GetUsuario", func(t *testing.T) {
		user, err := queries.GetUsuario(ctx, userID)
		if err != nil {
			t.Fatalf("Error al buscar usuario %v", err)
		}
		if user.Nombre != nombreUsuarioTest {
			t.Errorf("El nombre no es el que se ingreso, se esperaba %s y se obtuvo %s", nombreUsuarioTest, user.Nombre)
		}
	})

	t.Run("GetAllUsuarios", func(t *testing.T) {
		_, err := queries.GetAllUsuarios(ctx)
		if err != nil {
			t.Fatalf("Error al listar usuarios %v", err)
		}
	})

	nombreUpdate := "nombre update"
	emailUpdate := "update@test.com"
	passwordUpdate := "contra update"

	t.Run("UpdateUsuario", func(t *testing.T) {
		user, err := queries.UpdateUsuario(ctx, sqlc.UpdateUsuarioParams{
			ID:       userID,
			Nombre:   nombreUpdate,
			Email:    emailUpdate,
			Password: passwordUpdate,
		})
		if err != nil {
			t.Errorf("Error al actualizar usuario %v", err)
		}
		if user.Nombre != nombreUpdate || user.Email != emailUpdate || user.Password != passwordUpdate {
			t.Errorf("no se actualizaron los datos")
		}
	})

	t.Run("DeleteUsuario", func(t *testing.T) {
		err := queries.DeleteUsuario(ctx, userID)
		if err != nil {
			t.Errorf("Error al eliminar usuario %v", err)
		}
		_, err2 := queries.GetUsuario(ctx, userID)
		if err2 == nil {
			t.Errorf("No se elimino el usuario")
		}
	})

	// tests para producto
	var prodID int32
	nombreProductoTest := "producto test"
	descripcionTest := "descrpicion test"
	categoriaTest := "categoria test"
	colorTest := "color test"
	precioTest := int32(1234)

	t.Run("CreateProducto", func(t *testing.T) {
		prod, err := queries.CreateProducto(ctx, sqlc.CreateProductoParams{
			Nombre:      nombreProductoTest,
			Descripcion: descripcionTest,
			Categoria:   categoriaTest,
			Color:       colorTest,
			Precio:      precioTest,
		})
		if err != nil {
			t.Fatalf("Error al crear producto %v", err)
		}
		prodID = prod.ID
	})

	t.Run("GetProducto", func(t *testing.T) {
		prod, err := queries.GetProducto(ctx, prodID)
		if err != nil {
			t.Fatalf("Error al buscar producto %v", err)
		}
		if prod.Nombre != nombreProductoTest {
			t.Errorf("El nombre no es el que se ingreso, se esperaba %s y se obtuvo %s", nombreProductoTest, prod.Nombre)
		}
	})

	t.Run("GetAllProductos", func(t *testing.T) {
		_, err := queries.GetAllProductos(ctx)
		if err != nil {
			t.Fatalf("Error al listar productos %v", err)
		}
	})

	nombreProductoUpdate := "producto update"
	descripcionUpdate := "descrpicion update"
	categoriaUpdate := "categoria update"
	colorUpdate := "color update"
	precioUpdate := 5678

	t.Run("UpdateProducto", func(t *testing.T) {
		prod, err := queries.UpdateProducto(ctx, sqlc.UpdateProductoParams{
			ID:          prodID,
			Nombre:      nombreProductoUpdate,
			Descripcion: descripcionUpdate,
			Categoria:   categoriaUpdate,
			Color:       colorUpdate,
			Precio:      int32(precioUpdate),
		})
		if err != nil {
			t.Errorf("Error al actualizar producto %v", err)
		}
		if prod.Nombre != nombreProductoUpdate || prod.Descripcion != descripcionUpdate || prod.Categoria != categoriaUpdate || prod.Color != colorUpdate || prod.Precio != int32(precioUpdate) {
			t.Errorf("no se actualizaron los datos")
		}
	})

	t.Run("DeleteProducto", func(t *testing.T) {
		err := queries.DeleteProducto(ctx, prodID)
		if err != nil {
			t.Errorf("Error al eliminar producto %v", err)
		}
		_, err2 := queries.GetProducto(ctx, prodID)
		if err2 == nil {
			t.Errorf("No se elimino el producto")
		}
	})
}
