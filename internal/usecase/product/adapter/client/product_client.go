package client

import (
	"context"
	"fmt"
	"github.com/core-go/client"
	. "go-service/internal/usecase/product/domain"
	"net/http"
)

type ProductClient struct {
	Client *http.Client
	Url    string
	Config *client.LogConfig
	Log    func(context.Context, string, map[string]interface{})
}

type ResultInfo struct {
	Status  int64          `mapstructure:"status" json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status"`
	Errors  []ErrorMessage `mapstructure:"errors" json:"errors,omitempty" gorm:"column:errors" bson:"errors,omitempty" dynamodbav:"errors,omitempty" firestore:"errors,omitempty"`
	Message string         `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}
type ErrorMessage struct {
	Field   string `mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

func NewProductClient(config client.ClientConfig, log func(context.Context, string, map[string]interface{})) (*ProductClient, error) {
	c, _, conf, err := client.InitializeClient(config)
	if err != nil {
		return nil, err
	}
	return &ProductClient{Client: c, Url: config.Endpoint.Url, Config: conf, Log: log}, nil
}

func (c *ProductClient) Load(ctx context.Context, id string) (*Product, error) {
	url := c.Url + "/" + id
	var Product Product
	err := client.Get(ctx, c.Client, url, &Product, c.Config, c.Log)
	return &Product, err
}

func (c *ProductClient) Create(ctx context.Context, Product *Product) (int64, error) {
	var res ResultInfo
	err := client.Post(ctx, c.Client, c.Url, Product, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *ProductClient) Update(ctx context.Context, Product *Product) (int64, error) {
	url := c.Url + "/" + Product.Id
	var res ResultInfo
	err := client.Put(ctx, c.Client, url, Product, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *ProductClient) Patch(ctx context.Context, Product map[string]interface{}) (int64, error) {
	url := c.Url + "/" + fmt.Sprintf("%v", Product["id"])
	var res ResultInfo
	err := client.Patch(ctx, c.Client, url, Product, &res, c.Config, c.Log)
	return res.Status, err
}

func (c *ProductClient) Delete(ctx context.Context, id string) (int64, error) {
	url := c.Url + "/" + id
	var res int64
	err := client.Delete(ctx, c.Client, url, &res, c.Config, c.Log)
	return res, err
}
