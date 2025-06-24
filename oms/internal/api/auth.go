package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/omniful/go_commons/config"
)

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func TokenHandler(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	secret := config.GetString(c.Request.Context(), "jwt.secret")
	claims := jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not sign token"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: signed})
}
