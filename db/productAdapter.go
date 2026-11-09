package db

import (
	"context"

	sqlc "tpespecialweb/db/sqlc"
	"tpespecialweb/logic/imagenes"
	productcrud "tpespecialweb/logic/productCRUD"
)

type ProductAdapter struct {
	queries *sqlc.Queries
}

func NewProductAdapter(queries *sqlc.Queries) productcrud.ProductRepository {
	return &ProductAdapter{queries: queries}
}

func (a *ProductAdapter) Create(ctx context.Context, req productcrud.ProductoRequest) (productcrud.Producto, error) {
	// Mapeo desde la lógica hacia los parámetros específicos de sqlc
	params := sqlc.CreateProductoParams{
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Categoria:   req.Categoria,
		Color:       req.Color,
		Precio:      req.Precio,
	}

	res, err := a.queries.CreateProducto(ctx, params)
	if err != nil {
		return productcrud.Producto{}, err
	}

	return productcrud.Producto{
		ID:          res.ID,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Categoria:   res.Categoria,
		Color:       res.Color,
		Precio:      res.Precio,
		// Un producto recién creado todavía no tiene imágenes
		Imagenes: []imagenes.Imagen{},
	}, nil
}

func (a *ProductAdapter) Update(ctx context.Context, id int32, req productcrud.ProductoRequest) (productcrud.Producto, error) {
	params := sqlc.UpdateProductoParams{
		ID:          id,
		Nombre:      req.Nombre,
		Descripcion: req.Descripcion,
		Categoria:   req.Categoria,
		Color:       req.Color,
		Precio:      req.Precio,
	}

	res, err := a.queries.UpdateProducto(ctx, params)
	if err != nil {
		return productcrud.Producto{}, err
	}

	imgs, err := a.imagenesDe(ctx, id)
	if err != nil {
		return productcrud.Producto{}, err
	}

	return productcrud.Producto{
		ID:          id,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Categoria:   res.Categoria,
		Color:       res.Color,
		Precio:      res.Precio,
		Imagenes:    imgs,
	}, nil
}

func (a *ProductAdapter) GetByID(ctx context.Context, id int32) (productcrud.Producto, error) {
	res, err := a.queries.GetProducto(ctx, id)
	if err != nil {
		return productcrud.Producto{}, err
	}

	imgs, err := a.imagenesDe(ctx, id)
	if err != nil {
		return productcrud.Producto{}, err
	}

	return productcrud.Producto{
		ID:          id,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Categoria:   res.Categoria,
		Color:       res.Color,
		Precio:      res.Precio,
		Imagenes:    imgs,
	}, nil
}

func (a *ProductAdapter) GetAll(ctx context.Context) ([]productcrud.Producto, error) {
	list, err := a.queries.GetAllProductos(ctx)
	if err != nil {
		return nil, err
	}

	// Trae todas las imágenes en una sola consulta y las agrupa por producto,
	// en vez de hacer una consulta por cada producto (problema "N+1")
	todas, err := a.queries.GetAllProductoImagenes(ctx)
	if err != nil {
		return nil, err
	}
	porProducto := make(map[int32][]imagenes.Imagen)
	for _, img := range todas {
		porProducto[img.ProductoID] = append(porProducto[img.ProductoID], imagenes.Imagen{ID: img.ID, URL: img.Url})
	}

	// make(..., 0) para que una lista vacía se serialice como [] y no como null
	productos := make([]productcrud.Producto, 0, len(list))
	for _, res := range list {
		imgs := porProducto[res.ID]
		if imgs == nil {
			imgs = []imagenes.Imagen{}
		}
		productos = append(productos, productcrud.Producto{
			ID:          res.ID,
			Nombre:      res.Nombre,
			Descripcion: res.Descripcion,
			Categoria:   res.Categoria,
			Color:       res.Color,
			Precio:      res.Precio,
			Imagenes:    imgs,
		})
	}

	return productos, nil
}

func (a *ProductAdapter) Delete(ctx context.Context, id int32) error {
	return a.queries.DeleteProducto(ctx, id)
}

// imagenesDe devuelve las imágenes de un producto ([] si no tiene)
func (a *ProductAdapter) imagenesDe(ctx context.Context, productoID int32) ([]imagenes.Imagen, error) {
	list, err := a.queries.GetImagenesByProducto(ctx, productoID)
	if err != nil {
		return nil, err
	}
	imgs := make([]imagenes.Imagen, 0, len(list))
	for _, img := range list {
		imgs = append(imgs, imagenes.Imagen{ID: img.ID, URL: img.Url})
	}
	return imgs, nil
}

func (a *ProductAdapter) AddImagen(ctx context.Context, productoID int32, url string) (imagenes.Imagen, error) {
	res, err := a.queries.CreateProductoImagen(ctx, sqlc.CreateProductoImagenParams{
		ProductoID: productoID,
		Url:        url,
	})
	if err != nil {
		return imagenes.Imagen{}, err
	}
	return imagenes.Imagen{ID: res.ID, URL: res.Url}, nil
}

func (a *ProductAdapter) GetImagen(ctx context.Context, productoID int32, imagenID int32) (imagenes.Imagen, error) {
	res, err := a.queries.GetProductoImagen(ctx, sqlc.GetProductoImagenParams{
		ID:         imagenID,
		ProductoID: productoID,
	})
	if err != nil {
		return imagenes.Imagen{}, err
	}
	return imagenes.Imagen{ID: res.ID, URL: res.Url}, nil
}

func (a *ProductAdapter) DeleteImagen(ctx context.Context, productoID int32, imagenID int32) error {
	return a.queries.DeleteProductoImagen(ctx, sqlc.DeleteProductoImagenParams{
		ID:         imagenID,
		ProductoID: productoID,
	})
}
