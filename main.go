package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginAdminRequest struct {
	EmployeeID string `json:"employee_id"`
	Password   string `json:"password" binding:"required,gte=6,lte=24"`
}

type CreateUserRequest struct {
	Mobile    string `json:"mobile" binding:"required,len=10"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Age       int    `json:"age" binding:"required,gte=18,lte=65"`
	Gender    string `json:"gender" binding:"required,oneof=male female other"`
	Role      string `json:"role" binding:"required,oneof=employee"`
}

type User struct {
	EmployeeID string `json:"employee_id"`
	Mobile     string `json:"mobile"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	Role       string `json:"role"`
	Password   string `json:"password"`
}

var users []User
var sessionTokens []string

func generateEmployeeID() string {
	id, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("EMP%09d", id)
}

func generateSessionToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func contains(slice []string, str string) bool {
    for _, value := range slice {
        if value == str {
            return true
        }
    }
    return false
}

func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func loginAdmin(c *gin.Context) {
	var req LoginAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.EmployeeID == "admin" && req.Password == "Admin@123" {
		sessionToken := generateSessionToken()
		sessionTokens = append(sessionTokens, sessionToken)
		c.JSON(http.StatusOK, gin.H{"message": "Login successfully.", "session_token": sessionToken})
		return
	}


	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee id or password."})
}

func createUser(c *gin.Context) {
	if	authorizationToken := c.GetHeader("Authorization"); authorizationToken == "" || !contains(sessionTokens,authorizationToken){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session."})
		return
	}
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.ToLower(req.Role) != "employee" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role should be employee."})
		return
	}

	newUser := User{
		EmployeeID: generateEmployeeID(),
		Mobile:     req.Mobile,
		Email:      req.Email,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Age:        req.Age,
		Gender:     req.Gender,
		Role:       req.Role,
		Password:   generateRandomPassword(8),
	}
	users = append(users, newUser)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully.", "user": newUser})
}

func main() {
	r := gin.Default()
	r.POST("/login-admin", loginAdmin)
	r.POST("/create-user", createUser)
	r.Run(":8080")
}
