package views

import (
	"strconv"
	"strings"

	"tpespecialweb/logic/cueroCRUD"
	"tpespecialweb/logic/productCRUD"
)

// Vistas posibles del selector de la página principal (parámetro ?ver=...)
const (
	VistaCueros    = "cueros"
	VistaProductos = "productos"
)

// ElegirVista valida el parámetro ?ver= de la URL. Cualquier valor
// desconocido (o vacío) muestra los cueros.
func ElegirVista(ver string) string {
	if ver == VistaProductos {
		return VistaProductos
	}
	return VistaCueros
}

// IndexData es todo lo que necesita la página principal para renderizarse.
// Vista indica qué sección se muestra; solo hace falta cargar la lista de esa.
type IndexData struct {
	Vista        string
	Cueros       []cueroCRUD.Cuero
	Productos    []productCRUD.Producto
	FormCuero    CueroFormData
	FormProducto ProductoFormData
}

// CueroFormData guarda lo que el usuario escribió en el formulario de cuero.
// Los valores van como string para poder devolverlos tal cual si hubo un error.
type CueroFormData struct {
	Nombre      string
	Tipo        string
	Acabado     string
	Color       string
	EspesorMm   string
	Precio      string
	Descripcion string
	Error       string
}

// ProductoFormData guarda lo que el usuario escribió en el formulario de producto
type ProductoFormData struct {
	Nombre      string
	Categoria   string
	Color       string
	Precio      string
	Descripcion string
	Error       string
}

// Precio formatea un precio entero en pesos con separador de miles: 45000 -> "$ 45.000"
func Precio(precio int32) string {
	digitos := strconv.Itoa(int(precio))
	signo := ""
	if strings.HasPrefix(digitos, "-") {
		signo, digitos = "-", digitos[1:]
	}

	var b strings.Builder
	for i, d := range digitos {
		if i > 0 && (len(digitos)-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(d)
	}
	return signo + "$ " + b.String()
}

// Espesor formatea el espesor sin ceros de más: 2.5 -> "2,5 mm"
func Espesor(mm float64) string {
	return strings.Replace(strconv.FormatFloat(mm, 'f', -1, 64), ".", ",", 1) + " mm"
}

// unir junta los textos no vacíos con " · " (ej: "Vacuno · Graso")
func unir(partes ...string) string {
	noVacias := make([]string, 0, len(partes))
	for _, p := range partes {
		if strings.TrimSpace(p) != "" {
			noVacias = append(noVacias, p)
		}
	}
	return strings.Join(noVacias, " · ")
}
