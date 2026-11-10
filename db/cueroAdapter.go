package db

import (
	"context"

	sqlc "tpespecialweb/db/sqlc"
	cuerocrud "tpespecialweb/logic/cueroCRUD"
	"tpespecialweb/logic/imagenes"
)

type CueroAdapter struct {
	queries *sqlc.Queries
}

func NewCueroAdapter(queries *sqlc.Queries) cuerocrud.CueroRepository {
	return &CueroAdapter{queries: queries}
}

// toCuero mapea el struct generado por sqlc (más sus imágenes) al struct de la capa de lógica
func toCuero(c sqlc.Cuero, imgs []imagenes.Imagen) cuerocrud.Cuero {
	if imgs == nil {
		imgs = []imagenes.Imagen{} // para que se serialice como [] y no como null
	}
	return cuerocrud.Cuero{
		ID:          c.ID,
		Nombre:      c.Nombre,
		Tipo:        c.Tipo,
		Acabado:     c.Acabado,
		Color:       c.Color,
		EspesorMm:   c.EspesorMm,
		Descripcion: c.Descripcion,
		Precio:      c.Precio,
		Imagenes:    imgs,
	}
}

func toImagen(img sqlc.CueroImagen) imagenes.Imagen {
	return imagenes.Imagen{ID: img.ID, URL: img.Url}
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
	// Un cuero recién creado todavía no tiene imágenes
	return toCuero(res, nil), nil
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

	imgs, err := a.imagenesDe(ctx, id)
	if err != nil {
		return cuerocrud.Cuero{}, err
	}
	return toCuero(res, imgs), nil
}

func (a *CueroAdapter) GetByID(ctx context.Context, id int32) (cuerocrud.Cuero, error) {
	res, err := a.queries.GetCuero(ctx, id)
	if err != nil {
		return cuerocrud.Cuero{}, err
	}

	imgs, err := a.imagenesDe(ctx, id)
	if err != nil {
		return cuerocrud.Cuero{}, err
	}
	return toCuero(res, imgs), nil
}

func (a *CueroAdapter) GetAll(ctx context.Context) ([]cuerocrud.Cuero, error) {
	list, err := a.queries.GetAllCueros(ctx)
	if err != nil {
		return nil, err
	}

	// Trae todas las imágenes en una sola consulta y las agrupa por cuero,
	// en vez de hacer una consulta por cada cuero (problema "N+1")
	todas, err := a.queries.GetAllCueroImagenes(ctx)
	if err != nil {
		return nil, err
	}
	porCuero := make(map[int32][]imagenes.Imagen)
	for _, img := range todas {
		porCuero[img.CueroID] = append(porCuero[img.CueroID], toImagen(img))
	}

	cueros := make([]cuerocrud.Cuero, 0, len(list))
	for _, c := range list {
		cueros = append(cueros, toCuero(c, porCuero[c.ID]))
	}
	return cueros, nil
}

func (a *CueroAdapter) Delete(ctx context.Context, id int32) error {
	return a.queries.DeleteCuero(ctx, id)
}

// imagenesDe devuelve las imágenes de un cuero
func (a *CueroAdapter) imagenesDe(ctx context.Context, cueroID int32) ([]imagenes.Imagen, error) {
	list, err := a.queries.GetImagenesByCuero(ctx, cueroID)
	if err != nil {
		return nil, err
	}
	imgs := make([]imagenes.Imagen, 0, len(list))
	for _, img := range list {
		imgs = append(imgs, toImagen(img))
	}
	return imgs, nil
}

func (a *CueroAdapter) AddImagen(ctx context.Context, cueroID int32, url string) (imagenes.Imagen, error) {
	res, err := a.queries.CreateCueroImagen(ctx, sqlc.CreateCueroImagenParams{
		CueroID: cueroID,
		Url:     url,
	})
	if err != nil {
		return imagenes.Imagen{}, err
	}
	return toImagen(res), nil
}

func (a *CueroAdapter) GetImagen(ctx context.Context, cueroID int32, imagenID int32) (imagenes.Imagen, error) {
	res, err := a.queries.GetCueroImagen(ctx, sqlc.GetCueroImagenParams{
		ID:      imagenID,
		CueroID: cueroID,
	})
	if err != nil {
		return imagenes.Imagen{}, err
	}
	return toImagen(res), nil
}

func (a *CueroAdapter) DeleteImagen(ctx context.Context, cueroID int32, imagenID int32) error {
	return a.queries.DeleteCueroImagen(ctx, sqlc.DeleteCueroImagenParams{
		ID:      imagenID,
		CueroID: cueroID,
	})
}
