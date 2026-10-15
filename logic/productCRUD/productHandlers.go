package productCRUD

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	db "tpespecialweb/db/sqlc"
)

type ProductHandler struct {
	Queries *db.Queries
}

func NewProductHandler(q *db.Queries) *ProductHandler {
	return &ProductHandler{
		Queries: q,
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

	listaProductos, err := h.Queries.GetAllProductos(r.Context())

	if err != nil {
		http.Error(w, "Error al obtener los productos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(listaProductos)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	productoEncontrado, err := h.Queries.GetProducto(r.Context(), int32(id))

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Producto no encontrado", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productoEncontrado)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var nuevoProducto db.CreateProductoParams

	err := json.NewDecoder(r.Body).Decode(&nuevoProducto)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if nuevoProducto.Nombre == "" || nuevoProducto.Precio <= 0 {
		http.Error(w, "El nombre no puede estar vacío y el precio debe ser mayor a cero", http.StatusBadRequest)
		return
	}

	productoCreado, err := h.Queries.CreateProducto(r.Context(), nuevoProducto)

	if err != nil {
		http.Error(w, "Error al guardar el producto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(productoCreado)
}

func (h *ProductHandler) DeleteProductById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	err = h.Queries.DeleteProducto(r.Context(), int32(id))

	if err != nil {
		http.Error(w, "Error al eliminar el producto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (h *ProductHandler) UpdateProductById(w http.ResponseWriter, r *http.Request, idStr string) {

	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		http.Error(w, "El ID debe ser un número entero válido mayor a cero", http.StatusBadRequest)
		return
	}

	var pActualizar db.UpdateProductoParams

	if err := json.NewDecoder(r.Body).Decode(&pActualizar); err != nil {
		http.Error(w, "Formato JSON inválido", http.StatusBadRequest)
		return
	}

	pActualizar.ID = int32(id)

	if pActualizar.Nombre == "" || pActualizar.Precio <= 0 {
		http.Error(w, "El nombre no puede estar vacío y el precio debe ser mayor a cero", http.StatusBadRequest)
		return
	}

	productoActualizado, err := h.Queries.UpdateProducto(r.Context(), pActualizar)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Producto no encontrado", http.StatusNotFound) // 404
		} else {
			http.Error(w, "Error al actualizar el producto", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productoActualizado)
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")

	switch r.Method {
	case http.MethodGet:
		if idStr == "" {
			h.GetProducts(w, r)
		} else {
			h.GetProductById(w, r, idStr)
		}

	case http.MethodPost:
		h.CreateProduct(w, r)

	case http.MethodPut:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para actualizar", http.StatusBadRequest)
			return
		}
		h.UpdateProductById(w, r, idStr)

	case http.MethodDelete:
		if idStr == "" {
			http.Error(w, "Se requiere un ID para eliminar", http.StatusBadRequest)
			return
		}
		h.DeleteProductById(w, r, idStr)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
