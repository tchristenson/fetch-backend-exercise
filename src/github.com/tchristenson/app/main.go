package main

import (
	"fmt"
	"net/http"
	// "errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items`
	Total string `json:"total"`
}

type item struct {
	ID               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"Price"`
}

var receipts []receipt

func processReceipt(c *gin.Context) {
	uuid := uuid.New()
	var newReceipt receipt
	fmt.Println(uuid)

	if err := c.BindJSON(&newReceipt); err != nil {
		return
	}

	receipts = append(receipts, newReceipt)

	c.IndentedJSON(http.StatusCreated, newReceipt)
}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.Run("localhost:5000")
}
