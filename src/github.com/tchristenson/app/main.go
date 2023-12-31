package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type receipt struct {
	Id           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []struct {
		ShortDescription string `json:"shortDescription"`
		Price            string `json:"price"`
	} `json:"items"`
	Total string `json:"total"`
}

type item struct {
	Id               string `json:"id"`
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"Price"`
}

var receipts []receipt

func main() {
	router := gin.Default()
	router.POST("/receipts/process", processReceipt)
	router.GET("/receipts/:id/points", getReceiptPointsHandler)
	port := 5000
	router.Run(fmt.Sprintf("localhost:%d", port))
}

func processReceipt(c *gin.Context) {
	uuid := uuid.New()
	var newReceipt receipt
	newReceipt.Id = uuid.String()

	if err := c.BindJSON(&newReceipt); err != nil {
		c.String(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	// Check for missing or empty fields in the post body
	if strings.TrimSpace(newReceipt.Retailer) == "" ||
		strings.TrimSpace(newReceipt.PurchaseDate) == "" ||
		strings.TrimSpace(newReceipt.PurchaseTime) == "" ||
		len(newReceipt.Items) == 0 ||
		strings.TrimSpace(newReceipt.Total) == "" {
		c.String(http.StatusBadRequest, "The receipt is invalid")
		return
	}

	// Check for missing or empty item descriptions and prices
	for _, item := range newReceipt.Items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			c.String(http.StatusBadRequest, "The receipt is invalid")
			return
		}
		if strings.TrimSpace(item.Price) == "" {
			c.String(http.StatusBadRequest, "The receipt is invalid")
			return
		}
	}

	receipts = append(receipts, newReceipt)

	response := map[string]string{"id": newReceipt.Id}

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

func getReceiptPointsById(id string) (map[string]int64, error) {
	for _, receipt := range receipts {
		if receipt.Id == id {
			points := calculatePoints(receipt)
			return map[string]int64{"points": points}, nil
		}
	}

	return nil, fmt.Errorf("No receipt found for that id")
}

func calculatePoints(r receipt) int64 {
	points := int64(0)

	// One point for every alphanumeric character in the retailer name
	for _, char := range r.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points += 1
		}
	}

	// 50 points if the total is a round dollar amount with no cents, 25 points if the total is a multiple of 0.25
	totalFloat, err := strconv.ParseFloat(r.Total, 64)
	if err == nil {
		if int(totalFloat*100)%100 == 0 {
			points += 50
		}
		if int(totalFloat*100)%25 == 0 {
			points += 25
		}
	}

	// 5 points for every two items on the receipt
	numItems := len(r.Items)
	if numItems >= 2 {
		points += int64((numItems / 2) * 5)
	}

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer
	for _, item := range r.Items {
		trimmedItemDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedItemDescription)%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				points += int64(math.Ceil(priceFloat * 0.2))
			}
		}
	}

	// 6 points if the day in the purchase date is odd
	day := strings.Split(r.PurchaseDate, "-")[2]
	dayInt, err := strconv.Atoi(day)

	if err == nil {
		if dayInt%2 == 1 {
			points += 6
		}
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	hour := strings.Split(r.PurchaseTime, ":")[0]
	hourInt, err := strconv.Atoi(hour)

	if err == nil {
		minutes := strings.Split(r.PurchaseTime, ":")[1]
		minutesInt, err := strconv.Atoi(minutes)

		if err == nil {
			if (hourInt == 14 && minutesInt >= 1) || (hourInt == 15 && minutesInt <= 59) {
				points += 10
			}
		}
	}

	return points
}
