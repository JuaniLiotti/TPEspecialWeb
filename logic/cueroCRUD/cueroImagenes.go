package cueroCRUD

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"tpespecialweb/logic/imagenes"
)

// CarpetaImagenes devuelve la subcarpeta de uploads/ donde van las imágenes de un cuero
func CarpetaImagenes(cueroID int32) string {
	return fmt.Sprintf("cueros/%d", cueroID)
}

// RegistrarImagenes guarda en disco las imágenes ya validadas de un cuero y las
// asocia en la BD. La usan la API (POST /api/cueros/{id}/imagenes) y el
// formulario HTML de alta.
func RegistrarImagenes(ctx context.Context, repo CueroRepository, id int32, validadas []imagenes.Validada) ([]imagenes.Imagen, error) {
	urls, err := imagenes.Guardar(validadas, CarpetaImagenes(id))
	if err != nil {
		return nil, err
	}

	creadas := make([]imagenes.Imagen, 0, len(urls))
	for i, url := range urls {
		img, err := repo.AddImagen(ctx, id, url)
		if err != nil {
			// Si falla la BD, borra los archivos que todavía no quedaron registrados
			for _, pendiente := range urls[i:] {
				imagenes.Borrar(pendiente)
			}
			return nil, err
		}
		creadas = append(creadas, img)
	}
	return creadas, nil
}

// Eliminar borra un cuero y la carpeta con sus imágenes.
// Las filas de cuero_imagen se borran solas (ON DELETE CASCADE),
// pero los archivos en disco hay que borrarlos a mano.
func Eliminar(ctx context.Context, repo CueroRepository, id int32) error {
	if err := repo.Delete(ctx, id); err != nil {
		return err
	}
	if err := imagenes.BorrarCarpeta(CarpetaImagenes(id)); err != nil {
		log.Printf("No se pudieron borrar las imágenes del cuero %d: %v", id, err)
	}
	return nil
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

	headers, err := imagenes.Parsear(w, r)
	if err == nil && len(headers) == 0 {
		err = imagenes.ErrSinImagenes
	}
	var validadas []imagenes.Validada
	if err == nil {
		validadas, err = imagenes.Validar(headers)
	}
	if err != nil {
		if errors.Is(err, imagenes.ErrArchivoInvalido) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Error al leer las imágenes", http.StatusInternalServerError)
		}
		return
	}

	creadas, err := RegistrarImagenes(r.Context(), h.repo, id, validadas)
	if err != nil {
		log.Printf("Error al guardar imágenes del cuero %d: %v", id, err)
		http.Error(w, "Error al guardar las imágenes", http.StatusInternalServerError)
		return
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
