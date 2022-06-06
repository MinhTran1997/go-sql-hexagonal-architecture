package repository

import (
	"context"
	"database/sql"
	"errors"
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
	var productGeneral []ProductGeneral
	var productDetails []ProductDetails

	queryGeneral := fmt.Sprintf("select id, productName, description, price, status from products where id = %s limit 1", q.BuildParam(1))
	err := q.Query(ctx, r.DB, nil, &productGeneral, queryGeneral, id)
	if err != nil {
		return nil, err
	}

	queryDetails := fmt.Sprintf("select productID, supplier, storage, inStockAmount from product_details where productID = %s limit 1", q.BuildParam(1))
	err = q.Query(ctx, r.DB, nil, &productDetails, queryDetails, id)
	if err != nil {
		return nil, err
	}

	var product Product
	product.GeneralInfo = productGeneral[0]
	product.DetailInfo = productDetails[0]

	if &product != nil {
		return &product, nil
	}
	return nil, nil
}

func (r *ProductAdapter) Create(ctx context.Context, product *Product) (int64, error) {
	tx := GetTx(ctx)
	var rowsAffected int64

	if product.DetailInfo.InStockAmount > 0 {
		product.GeneralInfo.Status = "available"
	} else {
		product.GeneralInfo.Status = "not available"
	}

	queryGeneral, argsGeneral := q.BuildToInsert("products", product.GeneralInfo, q.BuildParam)
	_, errGeneral := tx.ExecContext(ctx, queryGeneral, argsGeneral...)
	if errGeneral != nil {
		return -1, errGeneral
	} else {
		rowsAffected++
	}

	if checkDetailReq := checkReqProductDetails(product.DetailInfo); checkDetailReq == nil {
		queryDetails, argsDetails := q.BuildToInsert("product_details", product.DetailInfo, q.BuildParam)
		_, errDetails := tx.ExecContext(ctx, queryDetails, argsDetails...)
		if errDetails != nil {
			return -1, errDetails
		} else {
			rowsAffected++
		}
	}

	return rowsAffected, nil
}

func (r *ProductAdapter) Update(ctx context.Context, product *Product) (int64, error) {
	tx := GetTx(ctx)
	var rowsAffected int64

	if product.DetailInfo.InStockAmount > 0 {
		product.GeneralInfo.Status = "available"
	} else {
		product.GeneralInfo.Status = "not available"
	}

	queryGeneral, argsGeneral := q.BuildToUpdate("products", product.GeneralInfo, q.BuildParam)
	_, err := tx.ExecContext(ctx, queryGeneral, argsGeneral...)
	if err != nil {
		return -1, err
	} else {
		rowsAffected++
	}

	if checkDetailReq := checkReqProductDetails(product.DetailInfo); checkDetailReq == nil {
		queryDetails, argsDetails := q.BuildToUpdate("product_details", product.DetailInfo, q.BuildParam)
		_, err1 := tx.ExecContext(ctx, queryDetails, argsDetails...)
		if err1 != nil {
			return -1, err1
		} else {
			rowsAffected++
		}
	}

	return rowsAffected, nil
}

func (r *ProductAdapter) Patch(ctx context.Context, product map[string]interface{}) (int64, error) {
	tx := GetTx(ctx)

	productType := reflect.TypeOf(ProductGeneral{})
	jsonColumnMap := q.MakeJsonColumnMap(productType)
	colMap := q.JSONToColumns(product, jsonColumnMap)
	keys, _ := q.FindPrimaryKeys(productType)

	query, args := q.BuildToPatch("products", colMap, keys, q.BuildParam)
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}

	return res.RowsAffected()
}

func (r *ProductAdapter) Delete(ctx context.Context, id string) (int64, error) {
	tx := GetTx(ctx)
	var rowsAffected int64

	queryDetails := fmt.Sprintf("delete from product_details where productId = %s", q.BuildParam(1))
	_, er1 := tx.ExecContext(ctx, queryDetails, id)
	if er1 != nil {
		return -1, er1
	} else {
		rowsAffected++
	}

	queryGeneral := fmt.Sprintf("delete from products where id = %s", q.BuildParam(1))
	_, er2 := tx.ExecContext(ctx, queryGeneral, id)
	if er2 != nil {
		return -1, er2
	} else {
		rowsAffected++
	}

	return rowsAffected, nil
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

func checkReqProductDetails(productDetails ProductDetails) error {
	count := 0
	v := reflect.ValueOf(productDetails)
	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		if reflect.ValueOf(values[i]).Kind() == reflect.String {
			if values[i] == "" {
				count++
			}
		}
		if reflect.ValueOf(values[i]).Kind() == reflect.Int {
			if values[i] == 0 {
				count++
			}
		}
	}

	if count == v.NumField() {
		return errors.New("product general request has no data")
	}
	return nil
}
