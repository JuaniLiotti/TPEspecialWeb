package cueroCRUD

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"tpespecialweb/logic/imagenes"
)

// carpetaImagenes devuelve la subcarpeta de uploads/ donde van las imágenes de un cuero
func carpetaImagenes(cueroID int32) string {
	return fmt.Sprintf("cueros/%d", cueroID)
}

// UploadImagenes maneja POST /api/cueros/{id}/imagenes
// Recibe multipart/form-data con uno o más archivos en el campo "imagenes".
func (h *CueroHandler) UploadImagenes(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// El cuero tiene que existir antes de guardar archivos
	if _, err := h.repo.GetByID(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Cuero no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	urls, err := imagenes.Recibir(w, r, carpetaImagenes(id))
	if err != nil {
		if errors.Is(err, imagenes.ErrArchivoInvalido) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			log.Printf("Error al guardar imágenes del cuero %d: %v", id, err)
			http.Error(w, "Error al guardar las imágenes", http.StatusInternalServerError)
		}
		return
	}

	creadas := make([]imagenes.Imagen, 0, len(urls))
	for i, url := range urls {
		img, err := h.repo.AddImagen(r.Context(), id, url)
		if err != nil {
			// Si falla la BD, borra los archivos que todavía no quedaron registrados
			for _, pendiente := range urls[i:] {
				imagenes.Borrar(pendiente)
			}
			http.Error(w, "Error al registrar las imágenes", http.StatusInternalServerError)
			return
		}
		creadas = append(creadas, img)
	}

	writeJSON(w, http.StatusCreated, creadas)
}

// DeleteImagen maneja DELETE /api/cueros/{id}/imagenes/{imgId}
func (h *CueroHandler) DeleteImagen(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	imgID, err := parseID(r.PathValue("imgId"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := h.repo.GetImagen(r.Context(), id, imgID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Imagen no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	if err := h.repo.DeleteImagen(r.Context(), id, imgID); err != nil {
		http.Error(w, "Error al eliminar la imagen", http.StatusInternalServerError)
		return
	}

	if err := imagenes.Borrar(img.URL); err != nil {
		log.Printf("No se pudo borrar el archivo %s: %v", img.URL, err)
	}

	w.WriteHeader(http.StatusNoContent)
}
