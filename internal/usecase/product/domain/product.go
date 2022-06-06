package domain

type ProductGeneral struct {
	Id          string `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"required,max=40" match:"equal"`
	ProductName string `json:"productName" gorm:"column:productName" bson:"productName" dynamodbav:"productName" firestore:"productName" avro:"productName" validate:"required,productName,max=100" match:"prefix"`
	Description string `json:"description" gorm:"column:description" bson:"description" dynamodbav:"description" firestore:"description" avro:"description" validate:"description,max=100" match:"prefix"`
	Price       string `json:"price" gorm:"column:price" bson:"price" dynamodbav:"price" firestore:"price" avro:"price" validate:"required,price,max=18"`
	Status      string `json:"status" gorm:"column:status" bson:"status" dynamodbav:"status" firestore:"status" avro:"status" validate:"description,max=100" match:"prefix"`
}

type ProductDetails struct {
	ProductID     string `json:"productID" gorm:"column:productID;primary_key" bson:"productID" dynamodbav:"productID" firestore:"productID" avro:"productID"`
	Supplier      string `json:"supplier" gorm:"column:supplier" bson:"supplier" dynamodbav:"supplier" firestore:"supplier" avro:"supplier"`
	Storage       string `json:"storage" gorm:"column:storage" bson:"storage" dynamodbav:"storage" firestore:"storage" avro:"storage"`
	InStockAmount int    `json:"inStockAmount" gorm:"column:inStockAmount" bson:"inStockAmount" dynamodbav:"inStockAmount" firestore:"inStockAmount" avro:"inStockAmount"`
}

type Product struct {
	GeneralInfo ProductGeneral
	DetailInfo  ProductDetails
}
