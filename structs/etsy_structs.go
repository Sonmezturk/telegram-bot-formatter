package structs

import "time"

type Customizations struct {
	Size            string `json:"Size"`
	Personalization string `json:"Personalization"`
}

type Item struct {
	Page           string         `json:"page"`
	Quantity       string         `json:"quantity"`
	SkuName        string         `json:"sku_name"`
	Color          string         `json:"color"`
	Customizations Customizations `json:"customizations"`
}

type Order struct {
	OrderDate string `json:"orderDate"`
	OrderID   string `json:"orderId"`
	Items     []Item `json:"items"`
	ShipTo    string `json:"shipTo"`
}

type FormatedOrder struct {
	Orders    []Order	`json:"orders"`
	Timestamp time.Time `json:"timestamp"` 
	UserName  string    `json:"userName"`  
	FileName  string	`json:"originalFilename"`  
}

type AggregatedItem struct {
	SkuName       string `bson:"_id.skuname"`
	Color         string `bson:"_id.color"`
	Size          string `bson:"_id.size"`
	TotalQuantity int    `bson:"totalQuantity"`
	Orders        []string `bson:"orders"`
}