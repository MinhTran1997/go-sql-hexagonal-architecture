package app

import (
	"context"
	. "github.com/core-go/service"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router, ctx context.Context, conf Config) error {
	app, err := NewApp(ctx, conf)
	if err != nil {
		return err
	}
	r.HandleFunc("/health", app.Health.Check).Methods(GET)

	product := "/products"
	r.HandleFunc(product+"/search", app.product.Search).Methods(GET, POST)
	r.HandleFunc(product+"/{id}", app.product.Load).Methods(GET)
	r.HandleFunc(product, app.product.Create).Methods(POST)
	r.HandleFunc(product+"/{id}", app.product.Update).Methods(PUT)
	r.HandleFunc(product+"/{id}", app.product.Patch).Methods(PATCH)
	r.HandleFunc(product+"/{id}", app.product.Delete).Methods(DELETE)

	return nil
}
