package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateExpense_BadRequest(t *testing.T) {
	router := gin.Default()
	router.POST("/api/v1/expenses", CreateExpense)

	reqBody := []byte(`{}`) // Invalid payload
	req, _ := http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
