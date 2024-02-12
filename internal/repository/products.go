package repository

//go:generate mockery --name=Products --inpackage --output=. --filename=products_mock.go --structname=ProductsMock

import (
	"errors"

	"golang.org/x/net/context"

	"github.com/riabininkf/go-project-template/internal/db"
	"github.com/riabininkf/go-project-template/internal/repository/entity"
)

func NewProducts(conn db.Queryer) Products {
	return &products{conn: conn}
}

type (
	Products interface {
		Add(ctx context.Context, product *entity.Product) error
		Get(ctx context.Context, id uint64) (*entity.Product, error)
		ExistsByName(ctx context.Context, name string) (bool, error)
	}

	products struct {
		conn db.Queryer
	}
)

func (r *products) Add(ctx context.Context, product *entity.Product) error {
	return r.conn.QueryRow(
		ctx,
		`
		INSERT INTO template.products (name, comment, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW()) RETURNING id
		`,
		product.Name,
		product.Comment,
	).Scan(&product.ID)
}

func (r *products) Get(ctx context.Context, id uint64) (*entity.Product, error) {
	var product entity.Product
	if err := r.conn.QueryRow(ctx, "SELECT id, name, comment FROM template.products WHERE id = $1", id).Scan(
		&product.ID,
		&product.Name,
		&product.Comment,
	); err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &product, nil
}

func (r *products) ExistsByName(ctx context.Context, name string) (bool, error) {
	var isExist bool
	if err := r.conn.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT TRUE FROM template.products WHERE name = $1)",
		name,
	).Scan(&isExist); err != nil {
		return false, err
	}

	return isExist, nil
}
