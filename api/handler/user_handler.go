package handler

import (
	"net/http"

	"db/service"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

// UserHandler represents the handler for user API
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,gt=0"`
}

// CreateUser handles the request to create a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), req.Name, req.Email, req.Age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Use jsoniter if available via middleware, otherwise fallback to standard c.JSON
	if jsonValue, exists := c.Get("json"); exists {
		if json, ok := jsonValue.(jsoniter.API); ok {
			data, err := json.Marshal(gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"age":   user.Age,
			})
			if err == nil {
				c.Data(http.StatusCreated, "application/json", data)
				return
			}
		}
	}

	// Fallback to standard Gin JSON marshaling
	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
}
