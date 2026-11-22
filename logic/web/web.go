// Package web tiene los handlers que devuelven páginas HTML renderizadas en el
// servidor con templ (TP5). Usan los mismos repositorios que la API JSON, pero
// en vez de responder JSON le pasan los datos a los componentes de views/.
package web

import (
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/a-h/templ"

	"tpespecialweb/logic/cueroCRUD"
	"tpespecialweb/logic/productCRUD"
	"tpespecialweb/views"
)

type WebHandler struct {
	productos productCRUD.ProductRepository
	cueros    cueroCRUD.CueroRepository
}

func NewWebHandler(productos productCRUD.ProductRepository, cueros cueroCRUD.CueroRepository) *WebHandler {
	return &WebHandler{productos: productos, cueros: cueros}
}

// Index maneja GET /. El parámetro ?ver=cueros|productos elige qué sección
// muestra el selector (por defecto, cueros). Es la única ruta que devuelve la
// página completa: los formularios (TP6) responden solo con fragmentos.
func (h *WebHandler) Index(w http.ResponseWriter, r *http.Request) {
	data := views.IndexData{Vista: views.ElegirVista(r.URL.Query().Get("ver"))}

	// Carga de la BD solo la lista de la vista elegida
	var err error
	if data.Vista == views.VistaProductos {
		data.Productos, err = h.productos.GetAll(r.Context())
	} else {
		data.Cueros, err = h.cueros.GetAll(r.Context())
	}
	if err != nil {
		log.Printf("Error al cargar la lista de %s: %v", data.Vista, err)
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}

	renderFragmentos(w, r, views.IndexPage(data))
}

// ---------------------------------------------------------------------
// Respuestas para HTMX (TP6)
// ---------------------------------------------------------------------

// renderFragmentos escribe uno o más componentes, uno detrás del otro, como
// respuesta HTML. Después de un alta se mandan dos: la lista (que va al
// hx-target) y el formulario con hx-swap-oob (que HTMX ubica por su id).
func renderFragmentos(w http.ResponseWriter, r *http.Request, componentes ...templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, c := range componentes {
		if err := c.Render(r.Context(), w); err != nil {
			log.Printf("Error al renderizar la respuesta: %v", err)
			return
		}
	}
}

// renderErrorFormulario devuelve solo el formulario (con el error y lo que
// había cargado el usuario) y le indica a HTMX que lo ponga en lugar del
// formulario de la página, en vez de en la lista que dice el hx-target:
//   - HX-Retarget: nuevo destino del swap (selector CSS)
//   - HX-Reswap:   nueva forma de hacer el swap
//
// Responde 200 a propósito: HTMX 1.9 no hace el swap de respuestas 4xx/5xx,
// y para cambiar eso habría que escribir JavaScript.
func renderErrorFormulario(w http.ResponseWriter, r *http.Request, idForm string, form templ.Component) {
	w.Header().Set("HX-Retarget", "#"+idForm)
	w.Header().Set("HX-Reswap", "outerHTML")
	renderFragmentos(w, r, form)
}

// ---------------------------------------------------------------------
// Utilidades para leer los formularios
// ---------------------------------------------------------------------

// parseID lee el {id} de rutas como DELETE /cueros/{id}
func parseID(idStr string) (int32, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil || id <= 0 {
		return 0, errors.New("el ID debe ser un número entero válido mayor a cero")
	}
	return int32(id), nil
}

// parsePrecio convierte el campo precio del formulario (un entero en pesos)
func parsePrecio(s string) (int32, error) {
	precio, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, errors.New("el precio debe ser un número entero")
	}
	return int32(precio), nil
}

// parseDecimal convierte un campo numérico con decimales. Acepta coma o punto
// y rechaza "NaN" e "Inf", que ParseFloat acepta pero no son medidas válidas.
func parseDecimal(s string) (float64, error) {
	n, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
	if err != nil || math.IsNaN(n) || math.IsInf(n, 0) {
		return 0, errors.New("debe ser un número")
	}
	return n, nil
}

// capitalizar pone en mayúscula la primera letra de un mensaje de error
// (los errores de Go se escriben en minúscula, pero en pantalla queda mejor así)
func capitalizar(s string) string {
	r, tam := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[tam:]
}
