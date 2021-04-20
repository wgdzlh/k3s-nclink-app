package controllers

import (
	"k3s-nclink-apps/data-source/entity"
	"k3s-nclink-apps/data-source/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController is for auth logic
type AuthController struct {
	userservice service.UserService
}

func (a AuthController) Login(c *gin.Context) {
	var loginInfo entity.User
	err := c.ShouldBindJSON(&loginInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.userservice.Find(&loginInfo)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found."})
		return
	}

	if user.Access != "rw" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User access limited."})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Password invalid."})
		return
	}

	token, err := a.userservice.GetJwtToken(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}