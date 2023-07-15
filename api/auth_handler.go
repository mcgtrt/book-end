package api

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mcgtrt/book-end/store"
	"github.com/mcgtrt/book-end/types"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore store.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

func (h *AuthHandler) HandleAuth(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return ErrBadRequest()
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		return ErrInvalidCredentials()
	}
	if !isPasswordValid(user.EncryptedPassword, authParams.Password) {
		return ErrInvalidCredentials()
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	expires := time.Now().Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID,
		"email":   user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return tokenString
}

func isPasswordValid(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func newAuthHandler(userStore store.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}
