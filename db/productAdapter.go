package db

import (
	"context"

	sqlc "tpespecialweb/db/sqlc"
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

	return productcrud.Producto{
		ID:          id,
		Nombre:      res.Nombre,
		Descripcion: res.Descripcion,
		Categoria:   res.Categoria,
		Color:       res.Color,
		Precio:      res.Precio,
	}, nil
}

func (a *ProductAdapter) GetByID(ctx context.Context, id int32) (productcrud.Producto, error) {
	res, err := a.queries.GetProducto(ctx, id)
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
	}, nil
}

func (a *ProductAdapter) GetAll(ctx context.Context) ([]productcrud.Producto, error) {
	list, err := a.queries.GetAllProductos(ctx)
	if err != nil {
		return nil, err
	}

	var productos []productcrud.Producto
	for _, res := range list {
		productos = append(productos, productcrud.Producto{
			Nombre:      res.Nombre,
			Descripcion: res.Descripcion,
			Categoria:   res.Categoria,
			Color:       res.Color,
			Precio:      res.Precio,
		})
	}

	return productos, nil
}

func (a *ProductAdapter) Delete(ctx context.Context, id int32) error {
	return a.queries.DeleteProducto(ctx, id)
}
