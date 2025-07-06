package controller

import (
	"expense-tracker/model"
	"expense-tracker/postgresql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

// CreateExpense godoc
// @Summary      Create an expense
// @Description  Create a new expense record
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        expense  body      model.Expense  true  "Expense data"
// @Success      201      {object}  model.Expense
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/expenses [post]
// @Security     BearerAuth
func CreateExpense(c *gin.Context) {

	var expense model.Expense

	if err := c.ShouldBindJSON(&expense); err != nil {
		log.Errorf("Unable to bind JSON, %v: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	expense.Id = uuid.New().String()
	expense.User_id = uuid.New().String()

	// Save to DB
	if err := postgresql.DB.Create(&expense).Error; err != nil {
		log.Errorf("Failed to create expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}
	log.WithFields(log.Fields{
		"user_id":    expense.User_id,
		"expense_id": expense.Id,
	}).Info("Created expense")
	log.Infof("Created Expense for user: %v with expense id: %v", expense.User_id, expense.Id)
	c.JSON(http.StatusCreated, gin.H{"expense": expense})

}

// GetExpenseById godoc
// @Summary      Get expense by ID
// @Description  Get a single expense by its ID
// @Tags         expenses
// @Produce      json
// @Param        id   path      string  true  "Expense ID"
// @Success      200  {object}  model.Expense
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/expenses/{id} [get]
// @Security     BearerAuth
func GetExpenseById(c *gin.Context) {

	// var expID struct {
	// 	expID string
	// }

	// err := c.ShouldBindJSON(&expID)
	// if err != nil {
	// 	log.Error("Unable to bind JSON, err: ", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
	// 	return
	// }

	var expense model.Expense

	id := c.Param("id")

	log.Infof("Fetching expense with ID: %s", id)

	if err := postgresql.DB.Where("Id=?", id).First(&expense).Error; err != nil {
		log.Errorf("Failed to fetch expense from DB: expense ID=%s, Error=%v", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to fetch expense from DB": id, "expense ID": err})

		return
	}
	log.Infof("Created Expense for user: %v with expense id: %v", expense.User_id, expense.Id)
	c.JSON(http.StatusOK, gin.H{"Created Expense for user": expense.User_id, "with expense": expense})

}

// UpdateExpense godoc
// @Summary      Update an expense
// @Description  Update an existing expense by ID
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        id      path      string         true  "Expense ID"
// @Param        expense body      model.Expense  true  "Expense data"
// @Success      200     {object}  model.Expense
// @Failure      400     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Router       /api/v1/expenses/{id} [put]
// @Security     BearerAuth
func UpdateExpense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Error("Missing expense ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing expense ID"})
		return
	}

	var updateData model.Expense
	if err := c.ShouldBindJSON(&updateData); err != nil {
		log.Errorf("Unable to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var expense model.Expense
	if err := postgresql.DB.Where("id = ?", id).First(&expense).Error; err != nil {
		log.Errorf("Expense not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
		return
	}

	// Update fields (except Id and User_id)
	expense.Amount = updateData.Amount
	expense.Currency = updateData.Currency
	expense.Category = updateData.Category
	expense.Description = updateData.Description
	expense.TimeStamp = updateData.TimeStamp

	if err := postgresql.DB.Save(&expense).Error; err != nil {
		log.Errorf("Failed to update expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expense"})
		return
	}

	log.Infof("Updated Expense with id: %v", id)
	c.JSON(http.StatusOK, gin.H{"message": "Expense updated", "expense": expense})
}

// DeleteExpense godoc
// @Summary      Delete an expense
// @Description  Delete an expense by ID
// @Tags         expenses
// @Produce      json
// @Param        id   path      string  true  "Expense ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /api/v1/expenses/{id} [delete]
// @Security     BearerAuth
func DeleteExpense(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		log.Error("Missing expense ID in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing expense ID"})
		return
	}

	if err := postgresql.DB.Delete(&model.Expense{}, "id = ?", id).Error; err != nil {
		log.Errorf("Failed to delete expense: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete expense"})
		return
	}

	log.Infof("Deleted Expense with id: %v", id)
	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
}

// ListExpensesWithFilters godoc
// @Summary      List expenses with filters
// @Description  List expenses with optional filters
// @Tags         expenses
// @Produce      json
// @Param        user_id   query     string  false  "User ID"
// @Param        category  query     string  false  "Category"
// @Param        currency  query     string  false  "Currency"
// @Param        from      query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to        query     string  false  "End date (YYYY-MM-DD)"
// @Success      200       {object}  []model.Expense
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/expenses [get]
// @Security     BearerAuth
func ListExpensesWithFilters(c *gin.Context) {
	var expenses []model.Expense
	query := postgresql.DB

	// Filters: user_id, category, currency, from, to
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if currency := c.Query("currency"); currency != "" {
		query = query.Where("currency = ?", currency)
	}
	if from := c.Query("from"); from != "" {
		query = query.Where("time_stamp >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		query = query.Where("time_stamp <= ?", to)
	}

	limit := 10
	offset := 0
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if o := c.Query("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&expenses).Error; err != nil {
		log.Errorf("Failed to list expenses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

// Summary godoc
// @Summary      Get expense summary
// @Description  Get summary of expenses by category, with optional filters
// @Tags         expenses
// @Produce      json
// @Param        user_id   query     string  false  "User ID"
// @Param        from      query     string  false  "Start date (YYYY-MM-DD)"
// @Param        to        query     string  false  "End date (YYYY-MM-DD)"
// @Success      200       {object}  map[string]float64
// @Failure      500       {object}  map[string]string
// @Router       /api/v1/expenses/summary [get]
// @Security     BearerAuth
func Summary(c *gin.Context) {
	type Result struct {
		Category string
		Total    float64
	}
	var results []Result
	query := postgresql.DB.Model(&model.Expense{})

	// Optional filters
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if from := c.Query("from"); from != "" {
		query = query.Where("time_stamp >= ?", from)
	}
	if to := c.Query("to"); to != "" {
		query = query.Where("time_stamp <= ?", to)
	}

	if err := query.Select("category, SUM(amount) as total").Group("category").Scan(&results).Error; err != nil {
		log.Errorf("Failed to summarize expenses: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to summarize expenses"})
		return
	}

	summary := make(map[string]float64)
	for _, r := range results {
		summary[r.Category] = r.Total
	}

	c.JSON(http.StatusOK, summary)
}
