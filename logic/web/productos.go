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

// CrearProducto maneja POST /productos (formulario "Nuevo producto").
// Mismo flujo que CrearCuero: validar todo, crear, guardar imágenes y
// redirect 303 a GET /?ver=productos (Post/Redirect/Get).
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
		h.errorProducto(w, r, http.StatusBadRequest, f, errParseo)
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
		h.errorProducto(w, r, http.StatusBadRequest, f, err)
		return
	}

	// 2. Validar las imágenes antes de crear el producto
	validadas, err := imagenes.Validar(headers)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, imagenes.ErrArchivoInvalido) {
			status = http.StatusBadRequest
		}
		h.errorProducto(w, r, status, f, err)
		return
	}

	// 3. Guardar el producto y después sus imágenes
	creado, err := h.productos.Create(r.Context(), req)
	if err != nil {
		log.Printf("Error al guardar el producto: %v", err)
		h.errorProducto(w, r, http.StatusInternalServerError, f, errors.New("no se pudo guardar el producto"))
		return
	}
	if len(validadas) > 0 {
		if _, err := productCRUD.RegistrarImagenes(r.Context(), h.productos, creado.ID, validadas); err != nil {
			log.Printf("Error al guardar imágenes del producto %d: %v", creado.ID, err)
			h.errorProducto(w, r, http.StatusInternalServerError, views.ProductoFormData{},
				fmt.Errorf("se guardó %q, pero no se pudieron guardar las imágenes", creado.Nombre))
			return
		}
	}

	// 4. Redirect: el navegador hace GET /?ver=productos (PRG)
	http.Redirect(w, r, "/?ver=productos", http.StatusSeeOther)
}

// EliminarProducto maneja POST /productos/{id}/eliminar
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

	http.Redirect(w, r, "/?ver=productos", http.StatusSeeOther)
}

// errorProducto vuelve a mostrar la sección de productos con el error y los valores cargados
func (h *WebHandler) errorProducto(w http.ResponseWriter, r *http.Request, status int, f views.ProductoFormData, err error) {
	f.Error = capitalizar(err.Error())
	h.renderIndex(w, r, status, views.IndexData{Vista: views.VistaProductos, FormProducto: f})
}
