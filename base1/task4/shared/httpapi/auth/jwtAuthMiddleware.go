package auth

import (
	"fmt"
	"net/http"
	"strings"
	"task4/shared/config"
	apperror "task4/shared/kernel/error"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authValue := c.GetHeader("Authorization")
		if authValue == "" || !strings.HasPrefix(authValue, "Bearer ") {
			c.JSON(http.StatusUnauthorized, apperror.InvalidTokenError.WriteDetail("Invalid Authorized Header"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authValue, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, apperror.InvalidTokenError.WriteDetail("Invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, apperror.InvalidTokenError.WriteDetail("Invalid token"))
			c.Abort()
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, apperror.InvalidTokenError.WriteDetail("Token expired"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Next()
	}
}
