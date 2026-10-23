package db

import (
	"context"

	sqlc "tpespecialweb/db/sqlc"
	usercrud "tpespecialweb/logic/userCRUD"
)

type UserAdapter struct {
	queries *sqlc.Queries
}

func NewUserAdapter(queries *sqlc.Queries) usercrud.UserRepository {
	return &UserAdapter{queries: queries}
}

func (a *UserAdapter) Create(ctx context.Context, req usercrud.UsuarioRequest) (usercrud.Usuario, error) {
	params := sqlc.CreateUsuarioParams{
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := a.queries.CreateUsuario(ctx, params)
	if err != nil {
		return usercrud.Usuario{}, err
	}

	return usercrud.Usuario{
		ID:     res.ID,
		Nombre: res.Nombre,
		Email:  res.Email,
	}, nil
}

func (a *UserAdapter) Update(ctx context.Context, id int32, req usercrud.UsuarioRequest) (usercrud.Usuario, error) {
	params := sqlc.UpdateUsuarioParams{
		ID:       id,
		Nombre:   req.Nombre,
		Email:    req.Email,
		Password: req.Password,
	}

	res, err := a.queries.UpdateUsuario(ctx, params)
	if err != nil {
		return usercrud.Usuario{}, err
	}

	return usercrud.Usuario{
		ID:       id,
		Nombre:   res.Nombre,
		Email:    res.Email,
		Password: res.Password,
	}, nil
}

func (a *UserAdapter) GetByID(ctx context.Context, id int32) (usercrud.Usuario, error) {
	res, err := a.queries.GetUsuario(ctx, id)
	if err != nil {
		return usercrud.Usuario{}, err
	}

	return usercrud.Usuario{
		ID:     id,
		Nombre: res.Nombre,
		Email:  res.Email,
	}, nil
}

func (a *UserAdapter) GetAll(ctx context.Context) ([]usercrud.Usuario, error) {
	list, err := a.queries.GetAllUsuarios(ctx) // O ListUsuarios según cómo se llame en sqlc
	if err != nil {
		return nil, err
	}

	var usuarios []usercrud.Usuario
	for _, res := range list {
		usuarios = append(usuarios, usercrud.Usuario{
			ID:     res.ID,
			Nombre: res.Nombre,
			Email:  res.Email,
		})
	}

	return usuarios, nil
}

func (a *UserAdapter) Delete(ctx context.Context, id int32) error {
	return a.queries.DeleteUsuario(ctx, id)
}
