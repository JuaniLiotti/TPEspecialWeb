package userCRUD

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

type UsuarioRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Usuario struct {
	ID       int32  `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	Create(ctx context.Context, req UsuarioRequest) (Usuario, error)
	GetByID(ctx context.Context, id int32) (Usuario, error)
	GetAll(ctx context.Context) ([]Usuario, error)
	Update(ctx context.Context, id int32, req UsuarioRequest) (Usuario, error)
	Delete(ctx context.Context, id int32) error
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	listaUsuarios, err := h.repo.GetAll(r.Context())

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

	usuarioEncontrado, err := h.repo.GetByID(r.Context(), int32(id))

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
	var nuevoUsuario UsuarioRequest

	err := json.NewDecoder(r.Body).Decode(&nuevoUsuario)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if nuevoUsuario.Nombre == "" || nuevoUsuario.Email == "" || nuevoUsuario.Password == "" {
		http.Error(w, "El nombre, email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	usuarioCreado, err := h.repo.Create(r.Context(), nuevoUsuario)

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

	err = h.repo.Delete(r.Context(), int32(id))

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

	var uActualizar UsuarioRequest

	if err := json.NewDecoder(r.Body).Decode(&uActualizar); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	if uActualizar.Nombre == "" || uActualizar.Email == "" || uActualizar.Password == "" {
		http.Error(w, "El nombre, email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	usuarioActualizado, err := h.repo.Update(r.Context(), int32(id), uActualizar)

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
