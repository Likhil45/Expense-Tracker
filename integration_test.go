package main

import (
	"bytes"
	"encoding/json"
	"expense-tracker/controller"
	"expense-tracker/model"
	"expense-tracker/postgresql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "user=user password=password dbname=expense_tracker_test host=localhost port=5432 sslmode=disable")
	os.Setenv("JWT_SECRET", "test_jwt_secret")
	postgresql.ConnectPostgres()
	setupTestDB()
	os.Exit(m.Run())
}

func setupTestDB() {
	postgresql.DB.AutoMigrate(&model.User{}, &model.Expense{})
	postgresql.DB.Exec("TRUNCATE TABLE users, expenses RESTART IDENTITY CASCADE;")
}

func TestExpenseIntegrationFlow(t *testing.T) {
	// Use test DB or your dev DB (make sure it's running)
	// postgresql.ConnectPostgres() // Moved to TestMain

	// Setup Gin router
	r := gin.Default()
	r.POST("/api/v1/signup", controller.SignUp)
	r.POST("/api/v1/login", controller.Login)
	r.POST("/api/v1/expenses", controller.CreateExpense)
	r.GET("/api/v1/expenses/:id", controller.GetExpenseById)

	// 1. Sign up a user
	signupBody := `{"user_name":"testuser","password":"testpass"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/signup", bytes.NewBuffer([]byte(signupBody)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var signupResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &signupResp)
	// userID := signupResp["user_id"].(string) // Not used, so removed to avoid compile error

	// 2. Login to get JWT
	loginBody := `{"user_name":"testuser","password":"testpass"}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer([]byte(loginBody)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var loginResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	token := loginResp["token"].(string)

	// 3. Create an expense
	expense := model.Expense{
		Amount:      42.0,
		Currency:    "USD",
		Category:    "Food",
		Description: "Lunch",
	}
	expenseJSON, _ := json.Marshal(expense)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/expenses", bytes.NewBuffer(expenseJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var expenseResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &expenseResp)
	expenseData := expenseResp["expense"].(map[string]interface{})
	expenseID := expenseData["Id"].(string)

	// 4. Get the expense by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/expenses/"+expenseID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var getResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &getResp)
	expenseObj, ok := getResp["expense"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'expense' in response, got: %v", getResp)
	}
	assert.Equal(t, expenseID, expenseObj["Id"])
}
