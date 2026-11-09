package db

import (
	"context"

	sqlc "tpespecialweb/db/sqlc"
	cuerocrud "tpespecialweb/logic/cueroCRUD"
)

type CueroAdapter struct {
	queries *sqlc.Queries
}

func NewCueroAdapter(queries *sqlc.Queries) cuerocrud.CueroRepository {
	return &CueroAdapter{queries: queries}
}

// toCuero mapea el struct generado por sqlc al struct de la capa de lógica
func toCuero(c sqlc.Cuero) cuerocrud.Cuero {
	return cuerocrud.Cuero{
		ID:          c.ID,
		Nombre:      c.Nombre,
		Tipo:        c.Tipo,
		Acabado:     c.Acabado,
		Color:       c.Color,
		EspesorMm:   c.EspesorMm,
		Descripcion: c.Descripcion,
		Precio:      c.Precio,
	}
}

func (a *CueroAdapter) Create(ctx context.Context, req cuerocrud.CueroRequest) (cuerocrud.Cuero, error) {
	res, err := a.queries.CreateCuero(ctx, sqlc.CreateCueroParams{
		Nombre:      req.Nombre,
		Tipo:        req.Tipo,
		Acabado:     req.Acabado,
		Color:       req.Color,
		EspesorMm:   req.EspesorMm,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	})
	if err != nil {
		return cuerocrud.Cuero{}, err
	}
	return toCuero(res), nil
}

func (a *CueroAdapter) Update(ctx context.Context, id int32, req cuerocrud.CueroRequest) (cuerocrud.Cuero, error) {
	res, err := a.queries.UpdateCuero(ctx, sqlc.UpdateCueroParams{
		ID:          id,
		Nombre:      req.Nombre,
		Tipo:        req.Tipo,
		Acabado:     req.Acabado,
		Color:       req.Color,
		EspesorMm:   req.EspesorMm,
		Descripcion: req.Descripcion,
		Precio:      req.Precio,
	})
	if err != nil {
		return cuerocrud.Cuero{}, err
	}
	return toCuero(res), nil
}

func (a *CueroAdapter) GetByID(ctx context.Context, id int32) (cuerocrud.Cuero, error) {
	res, err := a.queries.GetCuero(ctx, id)
	if err != nil {
		return cuerocrud.Cuero{}, err
	}
	return toCuero(res), nil
}

func (a *CueroAdapter) GetAll(ctx context.Context) ([]cuerocrud.Cuero, error) {
	list, err := a.queries.GetAllCueros(ctx)
	if err != nil {
		return nil, err
	}

	cueros := make([]cuerocrud.Cuero, 0, len(list))
	for _, c := range list {
		cueros = append(cueros, toCuero(c))
	}
	return cueros, nil
}

func (a *CueroAdapter) Delete(ctx context.Context, id int32) error {
	return a.queries.DeleteCuero(ctx, id)
}
