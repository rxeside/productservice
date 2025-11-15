package query

import (
	"context"
	"database/sql"

	"gitea.xscloud.ru/xscloud/golib/pkg/infrastructure/mysql"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	appmodel "productservice/pkg/product/application/model"
	"productservice/pkg/product/application/query"
)

func NewProductQueryService(client mysql.ClientContext) query.ProductQueryService {
	return &productQueryService{
		client: client,
	}
}

type productQueryService struct {
	client mysql.ClientContext
}

func (q *productQueryService) FindProduct(ctx context.Context, productID uuid.UUID) (*appmodel.Product, error) {
	productDTO := struct {
		ProductID uuid.UUID `db:"product_id"`
		Name      string    `db:"name"`
		Price     int64     `db:"price"`
	}{}

	err := q.client.GetContext(
		ctx,
		&productDTO,
		`SELECT product_id, name, price FROM product WHERE product_id = ?`,
		productID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.WithStack(err)
	}

	return &appmodel.Product{
		ProductID: productDTO.ProductID,
		Name:      productDTO.Name,
		Price:     productDTO.Price,
	}, nil
}
