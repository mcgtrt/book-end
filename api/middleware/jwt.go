package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var errUnauthorised = fmt.Errorf("unauthorised")

func JWTAuthenticate(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return errUnauthorised
	}
	claims := getClaims(token)
	if claims == nil {
		return errUnauthorised
	}
	expFloat := claims["expires"].(float64)
	expInt := int64(expFloat)

	if time.Now().Unix() > expInt {
		return fmt.Errorf("token expired")
	}

	return c.Next()
}

func getClaims(tokenString string) jwt.MapClaims {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnauthorised
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
