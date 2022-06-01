package port

import (
	"context"
	. "go-service/internal/usecase/product/domain"
)

type ProductRepository interface {
	Load(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) (int64, error)
	Update(ctx context.Context, product *Product) (int64, error)
	Patch(ctx context.Context, product map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}
