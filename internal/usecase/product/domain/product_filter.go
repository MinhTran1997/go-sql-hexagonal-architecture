package domain

import "github.com/core-go/search"

type ProductFilter struct {
	*search.Filter
	Id          string `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"id" avro:"id" validate:"required,max=40" match:"equal"`
	productName string `json:"productName" gorm:"column:productName" bson:"productName" dynamodbav:"productName" firestore:"productName" avro:"productName" validate:"required,productName,max=100" match:"prefix" q:"prefix"`
	description string `json:"description" gorm:"column:description" bson:"description" dynamodbav:"description" firestore:"description" avro:"description" validate:"description,max=100" match:"prefix" q:"prefix"`
	price       string `json:"price" gorm:"column:price" bson:"price" dynamodbav:"price" firestore:"price" avro:"price" validate:"required,price,max=18" q:"true"`
	status      string `json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status" avro:"status"`
}
