package relational

import (
	"context"
	"fmt"
	"super-descuentos/errs"
	"super-descuentos/model"
)

func (store *SQLStore) GetAuthors(ctx context.Context, offset, limit int) ([]model.User, error) {
	authors, err := store.Queries.GetAuthors(ctx, struct {
		Limit  int64
		Offset int64
	}{Limit: int64(limit), Offset: int64(offset)})

	if err != nil {
		fmt.Printf("hubo un problema al intentar obtener los autores: %v", err)
		return nil, fmt.Errorf("hubo un problema al intentar obtener los autores")
	}

	if len(authors) == 0 {
		return []model.User{}, errs.ErrAuthorNotFound
	}

	users := make([]model.User, 0, limit)

	for _, author := range authors {
		users = append(users, RepositoryAuthorToModel(author))
	}

	return users, nil
}
