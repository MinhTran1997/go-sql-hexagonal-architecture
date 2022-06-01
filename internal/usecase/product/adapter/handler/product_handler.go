package handler

import (
	"context"
	"encoding/json"
	"github.com/core-go/search"
	sv "github.com/core-go/service"
	"github.com/gorilla/mux"
	"net/http"
	"reflect"

	. "go-service/internal/usecase/product/domain"
	. "go-service/internal/usecase/product/service"
)

func NewProductHandler(find func(context.Context, interface{}, interface{}, int64, ...int64) (int64, string, error), service ProductService, logError func(context.Context, string)) *HttpProductHandler {
	filterType := reflect.TypeOf(ProductFilter{})
	modelType := reflect.TypeOf(Product{})
	searchHandler := search.NewSearchHandler(find, modelType, filterType, logError, nil)
	return &HttpProductHandler{service: service, SearchHandler: searchHandler}
}

type HttpProductHandler struct {
	service ProductService
	*search.SearchHandler
}

func (h *HttpProductHandler) Load(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	product, err := h.service.Load(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, product)
}
func (h *HttpProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var product Product
	er1 := json.NewDecoder(r.Body).Decode(&product)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}

	res, er2 := h.service.Create(r.Context(), &product)
	if er2 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusCreated, res)
}
func (h *HttpProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	var product Product
	er1 := json.NewDecoder(r.Body).Decode(&product)
	defer r.Body.Close()
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	if len(product.Id) == 0 {
		product.Id = id
	} else if id != product.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}

	res, er2 := h.service.Update(r.Context(), &product)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}
func (h *HttpProductHandler) Patch(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}

	var product Product
	productType := reflect.TypeOf(product)
	_, jsonMap, _ := sv.BuildMapField(productType)
	body, er1 := sv.BuildMapAndStruct(r, &product)
	if er1 != nil {
		http.Error(w, er1.Error(), http.StatusInternalServerError)
		return
	}
	if len(product.Id) == 0 {
		product.Id = id
	} else if id != product.Id {
		http.Error(w, "Id not match", http.StatusBadRequest)
		return
	}
	json, er2 := sv.BodyToJsonMap(r, product, body, []string{"id"}, jsonMap)
	if er2 != nil {
		http.Error(w, er2.Error(), http.StatusInternalServerError)
		return
	}

	res, er3 := h.service.Patch(r.Context(), json)
	if er3 != nil {
		http.Error(w, er3.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}
func (h *HttpProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		http.Error(w, "Id cannot be empty", http.StatusBadRequest)
		return
	}
	res, err := h.service.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	JSON(w, http.StatusOK, res)
}

func JSON(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(res)
}
