package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestProcessReceipt(t *testing.T) {
	r := SetUpRouter()
	r.POST("/receipts/process", processReceipt)
	companyId := uuid.New().String()
	newReceipt := receipt{
		Id:           companyId,
		Retailer:     "Walmart",
		PurchaseDate: "2023-03-08",
		PurchaseTime: "12:00:00",
		Items: []struct {
			ShortDescription string `json:"shortDescription"`
			Price            string `json:"price"`
		}{
			{
				ShortDescription: "Apples",
				Price:            "1.99",
			},
			{
				ShortDescription: "Oranges",
				Price:            "2.99",
			},
		},
		Total: "4.98",
	}

	jsonValue, _ := json.Marshal(newReceipt)
	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
