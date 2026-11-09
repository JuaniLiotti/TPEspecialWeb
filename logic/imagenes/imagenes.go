// Package imagenes maneja la subida, el guardado y el borrado de los archivos
// de imagen de productos y cueros. Los archivos se guardan en disco, dentro de
// DirBase, y en la base de datos solo se guarda su URL pública.
package imagenes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Imagen struct {
	ID  int32  `json:"id"`
	URL string `json:"url"`
}

const (
	// DirBase es la carpeta del disco donde se guardan las imágenes
	DirBase = "uploads"
	// URLBase es el prefijo con el que el servidor las publica (ver main.go)
	URLBase = "/uploads/"
	// MaxRequestBytes es el tamaño máximo de un request con imágenes (10 MB)
	MaxRequestBytes = 10 << 20
	// CampoFormulario es el nombre del campo multipart que trae los archivos
	CampoFormulario = "imagenes"
)

// Tipos aceptados y la extensión con la que se guarda cada uno
var tiposPermitidos = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// ErrArchivoInvalido indica un error del cliente (corresponde a un 400).
// Cualquier otro error devuelto por Recibir es un error del servidor (500).
var ErrArchivoInvalido = errors.New("archivo inválido")

// ErrSinImagenes indica que el request no trae ningún archivo en el campo "imagenes"
var ErrSinImagenes = fmt.Errorf("%w: no se envió ninguna imagen en el campo %q", ErrArchivoInvalido, CampoFormulario)

// archivoValidado es un archivo ya revisado, listo para guardarse
type archivoValidado struct {
	header    *multipart.FileHeader
	extension string
}

// Recibir lee los archivos del campo "imagenes" de un request multipart/form-data,
// valida que todos sean imágenes y los guarda en uploads/<carpeta>/.
// Devuelve la URL pública de cada archivo guardado.
// Primero valida todos los archivos y recién después guarda: si uno es inválido,
// no se guarda ninguno.
func Recibir(w http.ResponseWriter, r *http.Request, carpeta string) ([]string, error) {
	// Limita el tamaño del body para que no nos llenen la memoria o el disco
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBytes)
	if err := r.ParseMultipartForm(MaxRequestBytes); err != nil {
		return nil, fmt.Errorf("%w: el request debe ser multipart/form-data y pesar menos de 10 MB", ErrArchivoInvalido)
	}

	headers := r.MultipartForm.File[CampoFormulario]
	if len(headers) == 0 {
		return nil, ErrSinImagenes
	}

	// 1. Validar todos
	validados := make([]archivoValidado, 0, len(headers))
	for _, h := range headers {
		ext, err := detectarExtension(h)
		if err != nil {
			return nil, err
		}
		validados = append(validados, archivoValidado{header: h, extension: ext})
	}

	// 2. Guardar todos
	dir := filepath.Join(DirBase, carpeta)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(validados))
	for _, v := range validados {
		url, err := guardar(v, carpeta)
		if err != nil {
			// Si falla a mitad de camino, borra lo que ya se había guardado
			for _, u := range urls {
				Borrar(u)
			}
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// detectarExtension mira el contenido real del archivo (no la extensión que manda
// el cliente, que puede ser falsa) y devuelve la extensión con la que se guarda.
func detectarExtension(h *multipart.FileHeader) (string, error) {
	f, err := h.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// DetectContentType necesita como máximo los primeros 512 bytes
	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	tipo := http.DetectContentType(buf[:n])
	ext, ok := tiposPermitidos[tipo]
	if !ok {
		return "", fmt.Errorf("%w: %q no es una imagen JPG, PNG o WebP (tipo detectado: %s)", ErrArchivoInvalido, h.Filename, tipo)
	}
	return ext, nil
}

// guardar copia el archivo a disco con un nombre aleatorio generado por el servidor.
// Nunca se usa el nombre original: podría traer rutas como "../../main.go".
func guardar(v archivoValidado, carpeta string) (string, error) {
	nombre, err := nombreAleatorio()
	if err != nil {
		return "", err
	}
	nombre += v.extension

	origen, err := v.header.Open()
	if err != nil {
		return "", err
	}
	defer origen.Close()

	destino, err := os.Create(filepath.Join(DirBase, carpeta, nombre))
	if err != nil {
		return "", err
	}
	defer destino.Close()

	if _, err := io.Copy(destino, origen); err != nil {
		return "", err
	}
	return URLBase + carpeta + "/" + nombre, nil
}

func nombreAleatorio() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Borrar elimina del disco el archivo de una URL devuelta por Recibir
func Borrar(url string) error {
	ruta, ok := rutaEnDisco(url)
	if !ok {
		return fmt.Errorf("URL de imagen fuera de %s: %s", URLBase, url)
	}
	return os.Remove(ruta)
}

// BorrarCarpeta elimina todas las imágenes de una entidad (ej: "productos/3")
func BorrarCarpeta(carpeta string) error {
	return os.RemoveAll(filepath.Join(DirBase, carpeta))
}

// rutaEnDisco convierte "/uploads/productos/3/x.jpg" en "uploads/productos/3/x.jpg",
// verificando que la ruta no se escape de la carpeta de uploads.
func rutaEnDisco(url string) (string, bool) {
	if !strings.HasPrefix(url, URLBase) {
		return "", false
	}
	ruta := filepath.Join(DirBase, filepath.Clean("/"+strings.TrimPrefix(url, URLBase)))
	return ruta, strings.HasPrefix(ruta, DirBase+string(filepath.Separator))
}
