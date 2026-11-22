package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"tpespecialweb/logic/imagenes"
	"tpespecialweb/logic/productCRUD"
	"tpespecialweb/views"
)

// CrearProducto maneja POST /productos (formulario "Nuevo producto") por AJAX
// con HTMX. Mismo flujo que CrearCuero: validar todo, crear, guardar imágenes
// y responder con la lista (#lista-productos) más el formulario vacío fuera de
// banda; si hay un error, solo el formulario con el mensaje (#form-producto).
func (h *WebHandler) CrearProducto(w http.ResponseWriter, r *http.Request) {
	headers, errParseo := imagenes.Parsear(w, r)

	f := views.ProductoFormData{
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Categoria:   strings.TrimSpace(r.FormValue("categoria")),
		Color:       strings.TrimSpace(r.FormValue("color")),
		Precio:      strings.TrimSpace(r.FormValue("precio")),
		Descripcion: strings.TrimSpace(r.FormValue("descripcion")),
	}
	if errParseo != nil {
		errorProducto(w, r, f, errParseo)
		return
	}

	// 1. Convertir y validar los datos (mismas reglas que la API JSON)
	precio, err := parsePrecio(f.Precio)
	req := productCRUD.ProductoRequest{
		Nombre:      f.Nombre,
		Descripcion: f.Descripcion,
		Categoria:   f.Categoria,
		Color:       f.Color,
		Precio:      precio,
	}
	if err == nil {
		err = productCRUD.Validate(req)
	}
	if err != nil {
		errorProducto(w, r, f, err)
		return
	}

	// 2. Validar las imágenes antes de crear el producto
	validadas, err := imagenes.Validar(headers)
	if err != nil {
		if !errors.Is(err, imagenes.ErrArchivoInvalido) {
			log.Printf("Error al validar las imágenes del producto: %v", err)
		}
		errorProducto(w, r, f, err)
		return
	}

	// 3. Guardar el producto y después sus imágenes
	creado, err := h.productos.Create(r.Context(), req)
	if err != nil {
		log.Printf("Error al guardar el producto: %v", err)
		errorProducto(w, r, f, errors.New("no se pudo guardar el producto"))
		return
	}
	// El form que vuelve fuera de banda: vacío si salió todo bien
	formNuevo := views.ProductoFormData{}
	if len(validadas) > 0 {
		if _, err := productCRUD.RegistrarImagenes(r.Context(), h.productos, creado.ID, validadas); err != nil {
			log.Printf("Error al guardar imágenes del producto %d: %v", creado.ID, err)
			formNuevo.Error = capitalizar(fmt.Sprintf("se guardó %q, pero no se pudieron guardar las imágenes", creado.Nombre))
		}
	}

	// 4. Volver a consultar la lista y devolver los dos fragmentos
	productos, err := h.productos.GetAll(r.Context())
	if err != nil {
		log.Printf("Error al cargar la lista de productos: %v", err)
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}
	renderFragmentos(w, r, views.ProductoList(productos), views.ProductoFormOOB(formNuevo))
}

// EliminarProducto maneja DELETE /productos/{id}. Igual que EliminarCuero:
// respuesta vacía con 200 para que HTMX saque la tarjeta.
func (h *WebHandler) EliminarProducto(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := productCRUD.Eliminar(r.Context(), h.productos, id); err != nil {
		log.Printf("Error al eliminar el producto %d: %v", id, err)
		http.Error(w, "Error al eliminar el producto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// errorProducto devuelve el formulario de productos con el error y los valores
// cargados, para que HTMX lo ponga en #form-producto
func errorProducto(w http.ResponseWriter, r *http.Request, f views.ProductoFormData, err error) {
	f.Error = capitalizar(err.Error())
	renderErrorFormulario(w, r, "form-producto", views.ProductoForm(f))
}
