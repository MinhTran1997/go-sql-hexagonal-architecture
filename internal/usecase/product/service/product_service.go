package service

import (
	"context"
	"database/sql"
	"errors"
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

func NewProductService(db *sql.DB, repository ProductRepository) ProductService {
	return &productService{
		db:         db,
		repository: repository,
	}
}

type productService struct {
	db         *sql.DB
	repository ProductRepository
}

func (s *productService) Load(ctx context.Context, id string) (*Product, error) {
	return s.repository.Load(ctx, id)
}
func (s *productService) Create(ctx context.Context, product *Product) (int64, error) {
	err := checkProductGeneralReq(product.GeneralInfo)
	if err != nil {
		return -1, nil
	}
	tx, err := s.db.Begin()
	if err != nil {
		return -1, nil
	}
	ctx = context.WithValue(ctx, "tx", tx)
	res, err := s.repository.Create(ctx, product)
	if err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	if err = tx.Commit(); err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	return res, err
}
func (s *productService) Update(ctx context.Context, product *Product) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, nil
	}
	ctx = context.WithValue(ctx, "tx", tx)
	res, err := s.repository.Update(ctx, product)
	if err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	if err = tx.Commit(); err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	return res, err
}
func (s *productService) Patch(ctx context.Context, product map[string]interface{}) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, nil
	}
	ctx = context.WithValue(ctx, "tx", tx)
	res, err := s.repository.Patch(ctx, product)
	if err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	if err = tx.Commit(); err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	return res, err
}
func (s *productService) Delete(ctx context.Context, id string) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return -1, nil
	}
	ctx = context.WithValue(ctx, "tx", tx)
	res, err := s.repository.Delete(ctx, id)
	if err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	if err = tx.Commit(); err != nil {
		er2 := tx.Rollback()
		if er2 != nil {
			return -1, er2
		}
	}
	return res, err
}

func checkProductGeneralReq(productGeneral ProductGeneral) error {
	if productGeneral.Id == "" {
		return errors.New("product request has no productId")
	}

	//count := 0
	//v := reflect.ValueOf(productGeneral)
	//values := make([]interface{}, v.NumField())
	//
	//for i := 0; i < v.NumField(); i++ {
	//	values[i] = v.Field(i).Interface()
	//	if values[i] == "" {
	//		count++
	//	}
	//}
	//
	//if count == v.NumField() {
	//	return errors.New("product general request has no data")
	//}

	return nil
}
