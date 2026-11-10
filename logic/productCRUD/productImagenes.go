package productCRUD

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"tpespecialweb/logic/imagenes"
)

// CarpetaImagenes devuelve la subcarpeta de uploads/ donde van las imágenes de un producto
func CarpetaImagenes(productoID int32) string {
	return fmt.Sprintf("productos/%d", productoID)
}

// pathID lee un parámetro numérico de la ruta, ej: {id} en /api/products/{id}/imagenes
func pathID(r *http.Request, nombre string) (int32, bool) {
	id, err := strconv.Atoi(r.PathValue(nombre))
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
}

// RegistrarImagenes guarda en disco las imágenes ya validadas de un producto y
// las asocia en la BD. La usan la API (POST /api/products/{id}/imagenes) y el
// formulario HTML de alta.
func RegistrarImagenes(ctx context.Context, repo ProductRepository, id int32, validadas []imagenes.Validada) ([]imagenes.Imagen, error) {
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

// Eliminar borra un producto y la carpeta con sus imágenes.
// Las filas de producto_imagen se borran solas (ON DELETE CASCADE),
// pero los archivos en disco hay que borrarlos a mano.
func Eliminar(ctx context.Context, repo ProductRepository, id int32) error {
	if err := repo.Delete(ctx, id); err != nil {
		return err
	}
	if err := imagenes.BorrarCarpeta(CarpetaImagenes(id)); err != nil {
		log.Printf("No se pudieron borrar las imágenes del producto %d: %v", id, err)
	}
	return nil
}

// UploadImagenes maneja POST /api/products/{id}/imagenes
// Recibe multipart/form-data con uno o más archivos en el campo "imagenes".
func (h *ProductHandler) UploadImagenes(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	if !ok {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	// El producto tiene que existir antes de guardar archivos
	if _, err := h.repo.GetByID(r.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Producto no encontrado", http.StatusNotFound)
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
		log.Printf("Error al guardar imágenes del producto %d: %v", id, err)
		http.Error(w, "Error al guardar las imágenes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(creadas)
}

// DeleteImagen maneja DELETE /api/products/{id}/imagenes/{imgId}
func (h *ProductHandler) DeleteImagen(w http.ResponseWriter, r *http.Request) {
	id, ok := pathID(r, "id")
	imgID, ok2 := pathID(r, "imgId")
	if !ok || !ok2 {
		http.Error(w, "Los IDs deben ser números enteros válidos mayores a cero", http.StatusBadRequest)
		return
	}

	img, err := h.repo.GetImagen(r.Context(), id, imgID)
	if err != nil {
		if err == sql.ErrNoRows {
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

	w.WriteHeader(http.StatusNoContent) // 204
}
