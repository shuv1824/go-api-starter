package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shuv1824/go-api-starter/internal/common/errors"
	"github.com/shuv1824/go-api-starter/internal/domains/user/core"
)

type Handler struct {
	userService core.Service
}

func NewHandler(userService core.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req core.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.Register(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req core.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.userService.Login(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) handleError(c *gin.Context, err error) {
	switch err {
	case errors.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
	case errors.ErrEmailExists:
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
	case errors.ErrInvalidPassword:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	case errors.ErrUnauthorized:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	}
}
