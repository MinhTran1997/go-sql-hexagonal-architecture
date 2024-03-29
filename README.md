# go-sql-layer-architecture-sample

#### To run the application
```shell
go run main.go
```

## Architecture
### Simple Layer Architecture
![Layer Architecture](https://camo.githubproductcontent.com/d9b21eb50ef70dcaebf5a874559608f475e22c799bc66fcf99fb01f08576540f/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a4a4459546c4b3030796730496c556a5a392d737037512e706e67)

### Layer Architecture with full features
![Layer Architecture with standard features: config, health check, logging, middleware log tracing](https://camo.githubproductcontent.com/aa7b739a4692eaf2b363cf9caf8b021c60082c77c98d3f8c96665b5cf4640628/68747470733a2f2f63646e2d696d616765732d312e6d656469756d2e636f6d2f6d61782f3830302f312a6d79556b504969343265593477455f494446526176412e706e67)
#### [core-go/search](https://github.com/core-go/search)
- Build the search model at http handler
- Build dynamic SQL for search
  - Build SQL for paging by page index (page) and page size (limit)
  - Build SQL to count total of records
### Search products: Support both GET and POST 
#### POST /products/search
##### *Request:* POST /products/search
In the below sample, search products with these criteria:
- get products of page "1", with page size "20"
- description="tony": get products with description starting with "tony"
- status between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by price ascending, id descending
```json
{
    "page": 1,
    "limit": 20,
    "sort": "price,-id",
    "description": "tony",
    "status": {
        "min": "1953-11-16T00:00:00+07:00",
        "max": "1976-11-16T00:00:00+07:00"
    }
}
```
##### GET /products/search?page=1&limit=2&description=tony&status.min=1953-11-16T00:00:00+07:00&status.max=1976-11-16T00:00:00+07:00&sort=price,-id
In this sample, search products with these criteria:
- get products of page "1", with page size "20"
- description="tony": get products with description starting with "tony"
- status between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by price ascending, id descending

#### *Response:*
- total: total of products, which is used to calculate numbers of pages at client 
- list: list of products
```json
{
    "list": [
        {
            "id": "ironman",
            "productName": "tony.stark",
            "description": "tony.stark@gmail.com",
            "price": "0987654321",
            "status": "1963-03-24T17:00:00Z"
        }
    ],
    "total": 1
}
```

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource, refer to [service](https://github.com/core-go/service) and [mongo](https://github.com/core-go/mongo)  
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "mongo": {
            "status": "UP"
        }
    }
}
```

## API design for products
#### *Resource:* products

### Get all products
#### *Request:* GET /products
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "productName": "peter.parker",
        "description": "peter.parker@gmail.com",
        "price": "0987654321",
        "status": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "productName": "james.howlett",
        "description": "james.howlett@gmail.com",
        "price": "0987654321",
        "status": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one product by id
#### *Request:* GET /products/:id
```shell
GET /products/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "productName": "james.howlett",
    "description": "james.howlett@gmail.com",
    "price": "0987654321",
    "status": "1974-11-16T16:59:59.999Z"
}
```

### Create a new product
#### *Request:* POST /products 
```json
{
    "id": "wolverine",
    "productName": "james.howlett",
    "description": "james.howlett@gmail.com",
    "price": "0987654321",
    "status": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: duplicate key, -1: error
```json
1
```

### Update one product by id
#### *Request:* PUT /products/:id
```shell
PUT /products/wolverine
```
```json
{
    "productName": "james.howlett",
    "description": "james.howlett@gmail.com",
    "price": "0987654321",
    "status": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

### Patch one product by id
Perform a partial update of product. For example, if you want to update 2 fields: description and price, you can send the request body of below.
#### *Request:* PATCH /products/:id
```shell
PATCH /products/wolverine
```
```json
{
    "description": "james.howlett@gmail.com",
    "price": "0987654321"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

#### Problems for patch
If we pass a struct as a parameter, we cannot control what fields we need to update. So, we must pass a map as a parameter.
```go
type productservice interface {
    Update(ctx context.Context, product *product) (int64, error)
    Patch(ctx context.Context, product map[string]interface{}) (int64, error)
}
```
We must solve 2 problems:
1. At http handler layer, we must convert the product struct to map, with json format, and make sure the nested data types are passed correctly.
2. At repository layer, from json format, we must convert the json format to database format (in this case, we must convert to bson of Mongo)

#### Solutions for patch  
At http handler layer, we use [core-go/service](https://github.com/core-go/service), to convert the product struct to map, to make sure we just update the fields we need to update
```go
import server "github.com/core-go/service"

func (h *productHandler) Patch(w http.ResponseWriter, r *http.Request) {
    var product product
    productType := reflect.TypeOf(product)
    _, jsonMap := sv.BuildMapField(productType)
    body, _ := sv.BuildMapAndStruct(r, &product)
    json, er1 := sv.BodyToJson(r, product, body, ids, jsonMap, nil)

    result, er2 := h.service.Patch(r.Context(), json)
    if er2 != nil {
        http.Error(w, er2.Error(), http.StatusInternalServerError)
        return
    }
    respond(w, result)
}
```

### Delete a new product by id
#### *Request:* DELETE /products/:id
```shell
DELETE /products/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include HealthHandler, HealthChecker, SqlHealthChecker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): log and log middleware

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
    db, err := sql.Open(conf.Driver, conf.DataSourceName)
    if err != nil {
        return nil, err
    }

    sqlChecker := s.NewSqlHealthChecker(db)
    healthHandler := health.NewHealthHandler(sqlChecker)
```

To handler routing
```go
    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler.Check).Methods("GET")
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package main

import "github.com/core-go/config"

type Root struct {
    DB DatabaseConfig `mapstructure:"db"`
}

type DatabaseConfig struct {
    Driver         string `mapstructure:"driver"`
    DataSourceName string `mapstructure:"data_source_name"`
}

func main() {
    var conf Root
    err := config.Load(&conf, "configs/config")
    if err != nil {
        panic(err)
    }
}
```

### core-go/log *&* core-go/middleware
```go
import (
    "github.com/core-go/config"
    "github.com/core-go/log"
    m "github.com/core-go/middleware"
    "github.com/gorilla/mux"
)

func main() {
    var conf app.Root
    config.Load(&conf, "configs/config")

    r := mux.NewRouter()

    log.Initialize(conf.Log)
    r.Use(m.BuildContext)
    logger := m.NewStructuredLogger()
    r.Use(m.Logger(conf.MiddleWare, log.InfoFields, logger))
    r.Use(m.Recover(log.ErrorMsg))
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```
