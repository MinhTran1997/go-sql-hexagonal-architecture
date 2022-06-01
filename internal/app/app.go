package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	"github.com/core-go/search/query"
	q "github.com/core-go/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"

	"go-service/internal/usecase/product/adapter/handler"
	"go-service/internal/usecase/product/adapter/repository"
	. "go-service/internal/usecase/product/domain"
	// . "go-service/internal/client"
	. "go-service/internal/usecase/product/port"
	. "go-service/internal/usecase/product/service"
)

type ApplicationContext struct {
	Health  *health.Handler
	product ProductHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := q.OpenByConfig(conf.Sql)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg

	productType := reflect.TypeOf(Product{})
	productQueryBuilder := query.NewBuilder(db, "products", productType)
	productSearchBuilder, err := q.NewSearchBuilder(db, productType, productQueryBuilder.BuildQuery)
	if err != nil {
		return nil, err
	}
	/*
		client, _, _, err := client.InitializeClient(conf.Client)
		if err != nil {
			return nil, err
		}
		productRepository := NewproductClient(client, conf.Client.Endpoint.Url)*/
	productRepository := repository.NewProductAdapter(db)
	productService := NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productSearchBuilder.Search, productService, logError)

	sqlChecker := q.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health:  healthHandler,
		product: productHandler,
	}, nil
}
