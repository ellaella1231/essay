package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ai-essay-master-backend/internal/core/services"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type VerifyRequest struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

func (c *AuthController) Verify(ctx *gin.Context) {
	var req VerifyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := c.authService.VerifyOTP(ctx.Request.Context(), req.Phone, req.OTP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
