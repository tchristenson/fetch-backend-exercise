package main

import (
	"fmt"
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
)

type receipt struct{
	ID           string 	`json: "id"`
	Retailer 	 string 	`json: "retailer"`
	PurchaseDate string 	`json: "purchaseDate"`
	PurchaseTime string 	`json: "purchaseTime"`
	Items 		 [...]map   `json: "items"`
	Total 		 string 	`json: "total"`
}

type item struct{
	ID           	 string `json: "id"`
	ShortDescription string `json: "shortDescription"`
	Price 			 string `json: "Price"`
}

func main() {
	router := gin.Default()
}
