package userCRUD

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	db "tpespecialweb/db/sqlc"
)

type UserHandler struct {
	Queries *db.Queries
}

func NewUserHandler(q *db.Queries) *UserHandler {
	return &UserHandler{
		Queries: q,
	}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	listaUsuarios, err := h.Queries.GetAllUsuarios(r.Context())

	if err != nil {
		http.Error(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaUsuarios)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	usuarioEncontrado, err := h.Queries.GetUsuario(r.Context(), int32(id))

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioEncontrado)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var nuevoUsuario db.CreateUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&nuevoUsuario)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if nuevoUsuario.Nombre == "" || nuevoUsuario.Email == "" || nuevoUsuario.Password == "" {
		http.Error(w, "El nombre, email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	usuarioCreado, err := h.Queries.CreateUsuario(r.Context(), nuevoUsuario)

	if err != nil {
		http.Error(w, "Error al guardar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(usuarioCreado)
}

func (h *UserHandler) DeleteUserById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteUsuario(r.Context(), int32(id))

	if err != nil {
		http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	var uActualizar db.UpdateUsuarioParams

	if err := json.NewDecoder(r.Body).Decode(&uActualizar); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	uActualizar.ID = int32(id)

	if uActualizar.Nombre == "" || uActualizar.Email == "" || uActualizar.Password == "" {
		http.Error(w, "El nombre, email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	usuarioActualizado, err := h.Queries.UpdateUsuario(r.Context(), uActualizar)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarioActualizado)
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")

	switch r.Method {
	case http.MethodGet:
		if idStr == "" {
			h.GetUsers(w, r)
		} else {
			h.GetUserById(w, r, idStr)
		}

	case http.MethodPost:
		h.CreateUser(w, r)

	case http.MethodPut:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para actualizar", http.StatusBadRequest)
			return
		}
		h.UpdateUserById(w, r, idStr)

	case http.MethodDelete:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para eliminar", http.StatusBadRequest)
			return
		}
		h.DeleteUserById(w, r, idStr)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
