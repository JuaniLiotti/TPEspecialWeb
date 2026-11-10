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

	// tests para cuero
	var cueroID int32
	nombreCueroTest := "cuero test"
	tipoTest := "tipo test"
	acabadoTest := "acabado test"
	colorCueroTest := "color test"
	espesorTest := 2.5
	descripcionCueroTest := "descripcion test"
	precioCueroTest := int32(1000)

	t.Run("CreateCuero", func(t *testing.T) {
		cuero, err := queries.CreateCuero(ctx, sqlc.CreateCueroParams{
			Nombre:      nombreCueroTest,
			Tipo:        tipoTest,
			Acabado:     acabadoTest,
			Color:       colorCueroTest,
			EspesorMm:   espesorTest,
			Descripcion: descripcionCueroTest,
			Precio:      precioCueroTest,
		})
		if err != nil {
			t.Fatalf("Error al crear cuero %v", err)
		}
		cueroID = cuero.ID
	})

	t.Run("GetCuero", func(t *testing.T) {
		cuero, err := queries.GetCuero(ctx, cueroID)
		if err != nil {
			t.Fatalf("Error al buscar cuero %v", err)
		}
		if cuero.Nombre != nombreCueroTest || cuero.EspesorMm != espesorTest {
			t.Errorf("Los datos no son los que se ingresaron, se esperaba %s (%.1f mm) y se obtuvo %s (%.1f mm)", nombreCueroTest, espesorTest, cuero.Nombre, cuero.EspesorMm)
		}
	})

	t.Run("GetAllCueros", func(t *testing.T) {
		cueros, err := queries.GetAllCueros(ctx)
		if err != nil {
			t.Fatalf("Error al listar cueros %v", err)
		}
		encontrado := false
		for _, c := range cueros {
			if c.ID == cueroID {
				encontrado = true
			}
		}
		if !encontrado {
			t.Errorf("El cuero creado no aparece en la lista")
		}
	})

	nombreCueroUpdate := "cuero update"
	tipoUpdate := "tipo update"
	acabadoUpdate := "acabado update"
	colorCueroUpdate := "color update"
	espesorUpdate := 3.2
	descripcionCueroUpdate := "descripcion update"
	precioCueroUpdate := int32(2000)

	t.Run("UpdateCuero", func(t *testing.T) {
		cuero, err := queries.UpdateCuero(ctx, sqlc.UpdateCueroParams{
			ID:          cueroID,
			Nombre:      nombreCueroUpdate,
			Tipo:        tipoUpdate,
			Acabado:     acabadoUpdate,
			Color:       colorCueroUpdate,
			EspesorMm:   espesorUpdate,
			Descripcion: descripcionCueroUpdate,
			Precio:      precioCueroUpdate,
		})
		if err != nil {
			t.Errorf("Error al actualizar cuero %v", err)
		}
		if cuero.Nombre != nombreCueroUpdate || cuero.Tipo != tipoUpdate || cuero.Acabado != acabadoUpdate || cuero.Color != colorCueroUpdate || cuero.EspesorMm != espesorUpdate || cuero.Descripcion != descripcionCueroUpdate || cuero.Precio != precioCueroUpdate {
			t.Errorf("no se actualizaron los datos")
		}
	})

	t.Run("DeleteCuero", func(t *testing.T) {
		err := queries.DeleteCuero(ctx, cueroID)
		if err != nil {
			t.Errorf("Error al eliminar cuero %v", err)
		}
		_, err2 := queries.GetCuero(ctx, cueroID)
		if err2 == nil {
			t.Errorf("No se elimino el cuero")
		}
	})
}

func TestImagenes_CRUD(t *testing.T) {
	conn, err := sql.Open("postgres", "postgres://postgres:secreta@localhost:5432/italpiel_db?sslmode=disable")
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer conn.Close()

	queries := db.New(conn)
	ctx := context.Background()

	// tests para imágenes de producto
	prod, err := queries.CreateProducto(ctx, sqlc.CreateProductoParams{
		Nombre:      "producto con imagenes",
		Descripcion: "descripcion test",
		Categoria:   "categoria test",
		Color:       "color test",
		Precio:      1000,
	})
	if err != nil {
		t.Fatalf("Error al crear producto %v", err)
	}

	var prodImgID int32
	urlsProducto := []string{"/uploads/productos/test/1.jpg", "/uploads/productos/test/2.png"}

	t.Run("CreateProductoImagen", func(t *testing.T) {
		for i, url := range urlsProducto {
			img, err := queries.CreateProductoImagen(ctx, sqlc.CreateProductoImagenParams{ProductoID: prod.ID, Url: url})
			if err != nil {
				t.Fatalf("Error al crear imagen de producto %v", err)
			}
			if i == 0 {
				prodImgID = img.ID
			}
		}
	})

	t.Run("GetImagenesByProducto", func(t *testing.T) {
		imgs, err := queries.GetImagenesByProducto(ctx, prod.ID)
		if err != nil {
			t.Fatalf("Error al listar imagenes de producto %v", err)
		}
		if len(imgs) != 2 {
			t.Fatalf("Se esperaban 2 imagenes y se obtuvieron %d", len(imgs))
		}
		if imgs[0].Url != urlsProducto[0] {
			t.Errorf("La primera imagen (portada) deberia ser %s y es %s", urlsProducto[0], imgs[0].Url)
		}
	})

	t.Run("GetAllProductoImagenes", func(t *testing.T) {
		imgs, err := queries.GetAllProductoImagenes(ctx)
		if err != nil {
			t.Fatalf("Error al listar todas las imagenes de productos %v", err)
		}
		cantidad := 0
		for _, img := range imgs {
			if img.ProductoID == prod.ID {
				cantidad++
			}
		}
		if cantidad != 2 {
			t.Errorf("Se esperaban 2 imagenes del producto en la lista y hay %d", cantidad)
		}
	})

	t.Run("GetProductoImagen", func(t *testing.T) {
		img, err := queries.GetProductoImagen(ctx, sqlc.GetProductoImagenParams{ID: prodImgID, ProductoID: prod.ID})
		if err != nil {
			t.Fatalf("Error al buscar imagen de producto %v", err)
		}
		if img.Url != urlsProducto[0] {
			t.Errorf("Se esperaba %s y se obtuvo %s", urlsProducto[0], img.Url)
		}
	})

	t.Run("DeleteProductoImagen", func(t *testing.T) {
		err := queries.DeleteProductoImagen(ctx, sqlc.DeleteProductoImagenParams{ID: prodImgID, ProductoID: prod.ID})
		if err != nil {
			t.Errorf("Error al eliminar imagen de producto %v", err)
		}
		imgs, _ := queries.GetImagenesByProducto(ctx, prod.ID)
		if len(imgs) != 1 {
			t.Errorf("Deberia quedar 1 imagen y quedan %d", len(imgs))
		}
	})

	t.Run("DeleteProductoBorraSusImagenes", func(t *testing.T) {
		if err := queries.DeleteProducto(ctx, prod.ID); err != nil {
			t.Fatalf("Error al eliminar producto %v", err)
		}
		imgs, err := queries.GetImagenesByProducto(ctx, prod.ID)
		if err != nil {
			t.Fatalf("Error al listar imagenes de producto %v", err)
		}
		if len(imgs) != 0 {
			t.Errorf("ON DELETE CASCADE no borro las imagenes: quedan %d", len(imgs))
		}
	})

	// tests para imágenes de cuero
	cuero, err := queries.CreateCuero(ctx, sqlc.CreateCueroParams{
		Nombre:      "cuero con imagenes",
		Tipo:        "tipo test",
		Acabado:     "acabado test",
		Color:       "color test",
		EspesorMm:   2.0,
		Descripcion: "descripcion test",
		Precio:      1000,
	})
	if err != nil {
		t.Fatalf("Error al crear cuero %v", err)
	}

	var cueroImgID int32
	urlsCuero := []string{"/uploads/cueros/test/1.jpg", "/uploads/cueros/test/2.webp"}

	t.Run("CreateCueroImagen", func(t *testing.T) {
		for i, url := range urlsCuero {
			img, err := queries.CreateCueroImagen(ctx, sqlc.CreateCueroImagenParams{CueroID: cuero.ID, Url: url})
			if err != nil {
				t.Fatalf("Error al crear imagen de cuero %v", err)
			}
			if i == 0 {
				cueroImgID = img.ID
			}
		}
	})

	t.Run("GetImagenesByCuero", func(t *testing.T) {
		imgs, err := queries.GetImagenesByCuero(ctx, cuero.ID)
		if err != nil {
			t.Fatalf("Error al listar imagenes de cuero %v", err)
		}
		if len(imgs) != 2 {
			t.Fatalf("Se esperaban 2 imagenes y se obtuvieron %d", len(imgs))
		}
		if imgs[0].Url != urlsCuero[0] {
			t.Errorf("La primera imagen (portada) deberia ser %s y es %s", urlsCuero[0], imgs[0].Url)
		}
	})

	t.Run("GetCueroImagen", func(t *testing.T) {
		img, err := queries.GetCueroImagen(ctx, sqlc.GetCueroImagenParams{ID: cueroImgID, CueroID: cuero.ID})
		if err != nil {
			t.Fatalf("Error al buscar imagen de cuero %v", err)
		}
		if img.Url != urlsCuero[0] {
			t.Errorf("Se esperaba %s y se obtuvo %s", urlsCuero[0], img.Url)
		}
	})

	t.Run("DeleteCueroImagen", func(t *testing.T) {
		err := queries.DeleteCueroImagen(ctx, sqlc.DeleteCueroImagenParams{ID: cueroImgID, CueroID: cuero.ID})
		if err != nil {
			t.Errorf("Error al eliminar imagen de cuero %v", err)
		}
		imgs, _ := queries.GetImagenesByCuero(ctx, cuero.ID)
		if len(imgs) != 1 {
			t.Errorf("Deberia quedar 1 imagen y quedan %d", len(imgs))
		}
	})

	t.Run("DeleteCueroBorraSusImagenes", func(t *testing.T) {
		if err := queries.DeleteCuero(ctx, cuero.ID); err != nil {
			t.Fatalf("Error al eliminar cuero %v", err)
		}
		imgs, err := queries.GetImagenesByCuero(ctx, cuero.ID)
		if err != nil {
			t.Fatalf("Error al listar imagenes de cuero %v", err)
		}
		if len(imgs) != 0 {
			t.Errorf("ON DELETE CASCADE no borro las imagenes: quedan %d", len(imgs))
		}
	})
}
