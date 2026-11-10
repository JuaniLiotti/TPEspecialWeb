package cueroCRUD

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"tpespecialweb/logic/imagenes"
)

type Cuero struct {
	ID          int32   `json:"id"`
	Nombre      string  `json:"nombre"`
	Tipo        string  `json:"tipo"`
	Acabado     string  `json:"acabado"`
	Color       string  `json:"color"`
	EspesorMm   float64 `json:"espesor_mm"`
	Descripcion string  `json:"descripcion"`
	Precio      int32   `json:"precio"`
	// La primera imagen (la de menor ID) es la portada
	Imagenes []imagenes.Imagen `json:"imagenes"`
}

type CueroRequest struct {
	Nombre      string  `json:"nombre"`
	Tipo        string  `json:"tipo"`
	Acabado     string  `json:"acabado"`
	Color       string  `json:"color"`
	EspesorMm   float64 `json:"espesor_mm"`
	Descripcion string  `json:"descripcion"`
	Precio      int32   `json:"precio"`
}

type CueroRepository interface {
	Create(ctx context.Context, req CueroRequest) (Cuero, error)
	GetByID(ctx context.Context, id int32) (Cuero, error)
	GetAll(ctx context.Context) ([]Cuero, error)
	Update(ctx context.Context, id int32, req CueroRequest) (Cuero, error)
	Delete(ctx context.Context, id int32) error

	AddImagen(ctx context.Context, cueroID int32, url string) (imagenes.Imagen, error)
	GetImagen(ctx context.Context, cueroID int32, imagenID int32) (imagenes.Imagen, error)
	DeleteImagen(ctx context.Context, cueroID int32, imagenID int32) error
}

// Validate contiene las reglas de negocio de un cuero.
// La usan tanto la API JSON como los formularios HTML.
func Validate(req CueroRequest) error {
	if strings.TrimSpace(req.Nombre) == "" {
		return errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(req.Tipo) == "" {
		return errors.New("el tipo de cuero no puede estar vacío")
	}
	if req.EspesorMm <= 0 {
		return errors.New("el espesor debe ser mayor a cero")
	}
	if req.Precio <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	return nil
}

type CueroHandler struct {
	repo CueroRepository
}

func NewCueroHandler(repo CueroRepository) *CueroHandler {
	return &CueroHandler{repo: repo}
}

func parseID(idStr string) (int32, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("el ID debe ser un número entero válido mayor a cero")
	}
	return int32(id), nil
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *CueroHandler) GetCueros(w http.ResponseWriter, r *http.Request) {
	cueros, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los cueros", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, cueros)
}

func (h *CueroHandler) GetCueroById(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := parseID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cuero, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Cuero no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, cuero)
}

func (h *CueroHandler) CreateCuero(w http.ResponseWriter, r *http.Request) {
	var req CueroRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}
	if err := Validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	creado, err := h.repo.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "Error al guardar el cuero", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, creado)
}

func (h *CueroHandler) UpdateCueroById(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := parseID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req CueroRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}
	if err := Validate(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actualizado, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Cuero no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar el cuero", http.StatusInternalServerError)
		}
		return
	}
	writeJSON(w, http.StatusOK, actualizado)
}

func (h *CueroHandler) DeleteCueroById(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := parseID(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := Eliminar(r.Context(), h.repo, id); err != nil {
		http.Error(w, "Error al eliminar el cuero", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CueroHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/cueros/")

	switch r.Method {
	case http.MethodGet:
		if idStr == "" {
			h.GetCueros(w, r)
		} else {
			h.GetCueroById(w, r, idStr)
		}

	case http.MethodPost:
		h.CreateCuero(w, r)

	case http.MethodPut:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para actualizar", http.StatusBadRequest)
			return
		}
		h.UpdateCueroById(w, r, idStr)

	case http.MethodDelete:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para eliminar", http.StatusBadRequest)
			return
		}
		h.DeleteCueroById(w, r, idStr)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
