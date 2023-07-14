package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mcgtrt/book-end/api"
	"github.com/mcgtrt/book-end/store"
)

func JWTAuthenticate(userStore store.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return api.GenericResponseUnauthorised(c)
		}
		claims := getClaims(token)
		if claims == nil {
			return api.GenericResponseUnauthorised(c)
		}
		expFloat := claims["expires"].(float64)
		expInt := int64(expFloat)

		if time.Now().Unix() > expInt {
			return fmt.Errorf("token expired")
		}

		userID := claims["userID"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return api.GenericResponseUnauthorised(c)
		}
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func getClaims(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unauthorised")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return nil
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	return nil
}
