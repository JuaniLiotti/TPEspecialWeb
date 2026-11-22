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

// CrearCuero maneja POST /cueros (formulario "Nuevo cuero"), que llega por
// AJAX desde HTMX. Ya no hay redirect ni se recarga la página:
//   - si todo sale bien, responde con la lista actualizada (va al hx-target
//     #lista-cueros) y con el formulario vacío fuera de banda (lo limpia);
//   - si los datos son inválidos, responde solo con el formulario, con el error
//     y lo que el usuario había escrito (HX-Retarget lo manda a #form-cuero).
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
		errorCuero(w, r, f, errParseo)
		return
	}

	// 1. Convertir y validar los datos (mismas reglas que la API JSON)
	req, err := cueroRequest(f)
	if err == nil {
		err = cueroCRUD.Validate(req)
	}
	if err != nil {
		errorCuero(w, r, f, err)
		return
	}

	// 2. Validar las imágenes ANTES de crear el cuero, para no dejar un cuero
	// guardado cuando alguna foto se rechaza
	validadas, err := imagenes.Validar(headers)
	if err != nil {
		if !errors.Is(err, imagenes.ErrArchivoInvalido) {
			log.Printf("Error al validar las imágenes del cuero: %v", err)
		}
		errorCuero(w, r, f, err)
		return
	}

	// 3. Guardar el cuero y después sus imágenes
	creado, err := h.cueros.Create(r.Context(), req)
	if err != nil {
		log.Printf("Error al guardar el cuero: %v", err)
		errorCuero(w, r, f, errors.New("no se pudo guardar el cuero"))
		return
	}
	// El form que vuelve fuera de banda: vacío si salió todo bien
	formNuevo := views.CueroFormData{}
	if len(validadas) > 0 {
		if _, err := cueroCRUD.RegistrarImagenes(r.Context(), h.cueros, creado.ID, validadas); err != nil {
			log.Printf("Error al guardar imágenes del cuero %d: %v", creado.ID, err)
			// El cuero ya quedó creado: igual se actualiza la lista y el form vuelve
			// vacío (para que no lo carguen dos veces), pero con el aviso
			formNuevo.Error = capitalizar(fmt.Sprintf("se guardó %q, pero no se pudieron guardar las imágenes", creado.Nombre))
		}
	}

	// 4. Volver a consultar la lista y devolver los dos fragmentos
	cueros, err := h.cueros.GetAll(r.Context())
	if err != nil {
		log.Printf("Error al cargar la lista de cueros: %v", err)
		http.Error(w, "Error al obtener los datos", http.StatusInternalServerError)
		return
	}
	renderFragmentos(w, r, views.CueroList(cueros), views.CueroFormOOB(formNuevo))
}

// EliminarCuero maneja DELETE /cueros/{id} (botón "Eliminar" de cada tarjeta,
// que manda HTMX con hx-delete). Responde vacío con 200: HTMX reemplaza la
// tarjeta (hx-target="closest article") por nada y desaparece. No puede ser
// 204: con 204 HTMX 1.9 no hace el swap y la tarjeta quedaría en la página.
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

	w.WriteHeader(http.StatusOK)
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

// errorCuero devuelve el formulario de cueros con el mensaje de error y los
// valores que había cargado el usuario, para que HTMX lo ponga en #form-cuero
func errorCuero(w http.ResponseWriter, r *http.Request, f views.CueroFormData, err error) {
	f.Error = capitalizar(err.Error())
	renderErrorFormulario(w, r, "form-cuero", views.CueroForm(f))
}
