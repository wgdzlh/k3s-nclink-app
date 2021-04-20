package middlewares

import (
	"k3s-nclink-apps/data-source/service"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var userservice = service.UserService{}

// auth middleware
func AuthErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Authorization header is missing.",
			})
			return
		}

		temp := strings.Split(authHeader, "Bearer")
		if len(temp) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Wrong token."})
			return
		}

		tokenString := strings.TrimSpace(temp[1])
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return service.TokenKey, nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			name := claims["name"].(string)
			user, err := userservice.FindByName(name)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User not found."})
				return
			}
			if user.Access != "rw" {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User access limited."})
				return
			}
			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token invalid."})
		}
	}
}

// func GlobalErrorHandler(c *gin.Context) {
// 	c.Next()
// 	if len(c.Errors) > 0 {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": c.Errors})
// 	}
// }