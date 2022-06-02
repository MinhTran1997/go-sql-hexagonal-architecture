package repository

import (
	"context"
	"database/sql"
	"fmt"
	q "github.com/core-go/sql"
	. "go-service/internal/usecase/product/domain"
	"reflect"
)

func NewProductAdapter(db *sql.DB) *ProductAdapter {
	return &ProductAdapter{DB: db}
}

type ProductAdapter struct {
	DB *sql.DB
}

func (r *ProductAdapter) Load(ctx context.Context, id string) (*Product, error) {
	var products []Product
	query := fmt.Sprintf("select id, productName, description, price, status from products where id = %s limit 1", q.BuildParam(1))
	err := q.Query(ctx, r.DB, nil, &products, query, id)
	if err != nil {
		return nil, err
	}
	if len(products) > 0 {
		return &products[0], nil
	}
	return nil, nil
}

func (r *ProductAdapter) Create(ctx context.Context, product *Product) (int64, error) {
	query, args := q.BuildToInsert("products", product, q.BuildParam)

	tx := GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ProductAdapter) Update(ctx context.Context, product *Product) (int64, error) {
	query, args := q.BuildToUpdate("products", product, q.BuildParam)
	tx := GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ProductAdapter) Patch(ctx context.Context, product map[string]interface{}) (int64, error) {
	productType := reflect.TypeOf(Product{})
	jsonColumnMap := q.MakeJsonColumnMap(productType)
	colMap := q.JSONToColumns(product, jsonColumnMap)
	keys, _ := q.FindPrimaryKeys(productType)
	query, args := q.BuildToPatch("products", colMap, keys, q.BuildParam)

	tx := GetTx(ctx)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func (r *ProductAdapter) Delete(ctx context.Context, id string) (int64, error) {
	query := fmt.Sprintf("delete from products where id = %s", q.BuildParam(1))

	tx := GetTx(ctx)
	res, er1 := tx.ExecContext(ctx, query, id)
	if er1 != nil {
		return -1, er1
	}
	return res.RowsAffected()
}

func GetTx(ctx context.Context) *sql.Tx {
	txi := ctx.Value("tx")
	if txi != nil {
		txx, ok := txi.(*sql.Tx)
		if ok {
			return txx
		}
	}
	return nil
}
