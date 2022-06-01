package service

import (
	"context"

	. "go-service/internal/usecase/product/domain"
	. "go-service/internal/usecase/product/port"
)

type ProductService interface {
	Load(ctx context.Context, id string) (*Product, error)
	Create(ctx context.Context, product *Product) (int64, error)
	Update(ctx context.Context, product *Product) (int64, error)
	Patch(ctx context.Context, product map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id string) (int64, error)
}

func NewProductService(repository ProductRepository) ProductService {
	return &productService{repository: repository}
}

type productService struct {
	repository ProductRepository
}

func (s *productService) Load(ctx context.Context, id string) (*Product, error) {
	return s.repository.Load(ctx, id)
}
func (s *productService) Create(ctx context.Context, product *Product) (int64, error) {
	return s.repository.Create(ctx, product)
}
func (s *productService) Update(ctx context.Context, product *Product) (int64, error) {
	return s.repository.Update(ctx, product)
}
func (s *productService) Patch(ctx context.Context, product map[string]interface{}) (int64, error) {
	return s.repository.Patch(ctx, product)
}
func (s *productService) Delete(ctx context.Context, id string) (int64, error) {
	return s.repository.Delete(ctx, id)
}
