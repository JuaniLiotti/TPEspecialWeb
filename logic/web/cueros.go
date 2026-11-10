package web

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"tpespecialweb/logic/cueroCRUD"
	"tpespecialweb/logic/imagenes"
	"tpespecialweb/views"
)

// CrearCuero maneja POST /cueros (formulario "Nuevo cuero").
//
// Sigue el patrón Post/Redirect/Get: si todo sale bien responde con un
// redirect 303 a GET /?ver=cueros, así el navegador pide la página de nuevo
// (con la lista actualizada) y un F5 no vuelve a enviar el formulario.
// Si los datos son inválidos, vuelve a mostrar la página con el error y con
// lo que el usuario había escrito.
func (h *WebHandler) CrearCuero(w http.ResponseWriter, r *http.Request) {
	// Parsear limita el tamaño y lee el multipart/form-data (campos + archivos)
	headers, errParseo := imagenes.Parsear(w, r)

	f := views.CueroFormData{
		Nombre:      strings.TrimSpace(r.FormValue("nombre")),
		Tipo:        strings.TrimSpace(r.FormValue("tipo")),
		Acabado:     strings.TrimSpace(r.FormValue("acabado")),
		Color:       strings.TrimSpace(r.FormValue("color")),
		EspesorMm:   strings.TrimSpace(r.FormValue("espesor_mm")),
		Precio:      strings.TrimSpace(r.FormValue("precio")),
		Descripcion: strings.TrimSpace(r.FormValue("descripcion")),
	}
	if errParseo != nil {
		h.errorCuero(w, r, http.StatusBadRequest, f, errParseo)
		return
	}

	// 1. Convertir y validar los datos (mismas reglas que la API JSON)
	req, err := cueroRequest(f)
	if err == nil {
		err = cueroCRUD.Validate(req)
	}
	if err != nil {
		h.errorCuero(w, r, http.StatusBadRequest, f, err)
		return
	}

	// 2. Validar las imágenes ANTES de crear el cuero, para no dejar un cuero
	// guardado cuando alguna foto se rechaza
	validadas, err := imagenes.Validar(headers)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, imagenes.ErrArchivoInvalido) {
			status = http.StatusBadRequest
		}
		h.errorCuero(w, r, status, f, err)
		return
	}

	// 3. Guardar el cuero y después sus imágenes
	creado, err := h.cueros.Create(r.Context(), req)
	if err != nil {
		log.Printf("Error al guardar el cuero: %v", err)
		h.errorCuero(w, r, http.StatusInternalServerError, f, errors.New("no se pudo guardar el cuero"))
		return
	}
	if len(validadas) > 0 {
		if _, err := cueroCRUD.RegistrarImagenes(r.Context(), h.cueros, creado.ID, validadas); err != nil {
			log.Printf("Error al guardar imágenes del cuero %d: %v", creado.ID, err)
			// El cuero ya quedó creado: se vacía el form para que no lo carguen dos veces
			h.errorCuero(w, r, http.StatusInternalServerError, views.CueroFormData{},
				fmt.Errorf("se guardó %q, pero no se pudieron guardar las imágenes", creado.Nombre))
			return
		}
	}

	// 4. Redirect: el navegador hace GET /?ver=cueros (PRG)
	http.Redirect(w, r, "/?ver=cueros", http.StatusSeeOther)
}

// EliminarCuero maneja POST /cueros/{id}/eliminar (botón "Eliminar" de cada tarjeta).
// Los formularios HTML solo pueden mandar GET o POST, por eso no es un DELETE.
func (h *WebHandler) EliminarCuero(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := cueroCRUD.Eliminar(r.Context(), h.cueros, id); err != nil {
		log.Printf("Error al eliminar el cuero %d: %v", id, err)
		http.Error(w, "Error al eliminar el cuero", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/?ver=cueros", http.StatusSeeOther)
}

// cueroRequest convierte los textos del formulario al request que usa la lógica
func cueroRequest(f views.CueroFormData) (cueroCRUD.CueroRequest, error) {
	espesor, err := parseDecimal(f.EspesorMm)
	if err != nil {
		return cueroCRUD.CueroRequest{}, fmt.Errorf("el espesor %w", err)
	}
	precio, err := parsePrecio(f.Precio)
	if err != nil {
		return cueroCRUD.CueroRequest{}, err
	}
	return cueroCRUD.CueroRequest{
		Nombre:      f.Nombre,
		Tipo:        f.Tipo,
		Acabado:     f.Acabado,
		Color:       f.Color,
		EspesorMm:   espesor,
		Descripcion: f.Descripcion,
		Precio:      precio,
	}, nil
}

// errorCuero vuelve a mostrar la sección de cueros con el mensaje de error
// en el formulario y los valores que había cargado el usuario
func (h *WebHandler) errorCuero(w http.ResponseWriter, r *http.Request, status int, f views.CueroFormData, err error) {
	f.Error = capitalizar(err.Error())
	h.renderIndex(w, r, status, views.IndexData{Vista: views.VistaCueros, FormCuero: f})
}
