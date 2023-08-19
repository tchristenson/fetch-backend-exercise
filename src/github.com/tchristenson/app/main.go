package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"unicode"
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
	newReceipt.ID = uuid.String()

	if err := c.BindJSON(&newReceipt); err != nil {
		c.String(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	if len(newReceipt.Items) == 0 {
		c.String(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	receipts = append(receipts, newReceipt)
	// fmt.Println(receipts)

	response := map[string]string{"id": newReceipt.ID}

	c.IndentedJSON(http.StatusOK, response)
}

func getReceiptPointsHandler(c *gin.Context) {
	id := c.Param("id")
	points, err := getReceiptPointsById(id)

	if err != nil {
		c.String(http.StatusNotFound, "No receipt found for that id")
		return
	}

	c.IndentedJSON(http.StatusOK, points)
}

func getReceiptPointsById(id string) (map[string]int, error) {
	for _, receipt := range receipts {
		if receipt.ID == id {
			points := calculatePoints(receipt)
			return map[string]int{"points": points}, nil
		}
	}

	return nil, fmt.Errorf("No receipt found for that id")
}

func calculatePoints(r receipt) int {
	points := 0

	// Check alphanumeric characters
	for _, char := range r.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}

	// Check if the total is a round dollar amount
	totalFloat, err := strconv.ParseFloat(r.Total, 64)
	fmt.Println(totalFloat)
	if err == nil {
		if int(totalFloat*100)%100 == 0 {
			points += 50
		}
	}

	return points
}

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptPointsHandler)
	router.Run("localhost:5000")
}
