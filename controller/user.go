package controller

import (
	"expense-tracker/auth"
	"expense-tracker/model"
	"expense-tracker/postgresql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary      Create a new user (admin use)
// @Description  Create a new user with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  object  true  "User credentials"  example({"user_name":"admin","password":"secret"})
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/users [post]
// @Security     BearerAuth
func CreateUser(c *gin.Context) {
	var req struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := model.User{
		UserId:   uuid.New().String(),
		UserName: req.UserName,
		Password: string(hashed),
	}
	if err := postgresql.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user.UserName, "user_id": user.UserId})
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Register a new user with username and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  object  true  "User credentials"  example({"user_name":"alice","password":"mypassword"})
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/signup [post]
func SignUp(c *gin.Context) {
	var req struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := model.User{
		UserId:   uuid.New().String(),
		UserName: req.UserName,
		Password: string(hashed),
	}
	if err := postgresql.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign up"})
		return
	}
	token, err := auth.GenerateToken(user.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user.UserName, "user_id": user.UserId, "token": token})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  object  true  "User credentials"  example({"user_name":"alice","password":"mypassword"})
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/login [post]
func Login(c *gin.Context) {
	var req struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.UserName == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	var user model.User
	if err := postgresql.DB.Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	token, err := auth.GenerateToken(user.UserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user.UserName, "user_id": user.UserId, "token": token})
}
