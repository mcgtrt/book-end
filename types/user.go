package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	maxFirstNameLen = 64
	minLastNameLen  = 2
	maxLastNameLen  = 64
	minEmailLen     = 7
	maxEmailLen     = 128
	minPasswordLen  = 8
	maxPasswordLen  = 32
)

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(p.FirstName) < minFirstNameLen || len(p.FirstName) > maxFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name must be between %d - %d characters", minFirstNameLen, maxFirstNameLen)
	}
	if len(p.LastName) < minLastNameLen || len(p.LastName) > maxLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name must be between %d - %d characters", minLastNameLen, maxLastNameLen)
	}
	return errors
}

func (p UpdateUserParams) ToBSON() bson.M {
	b := bson.M{}
	if p.FirstName != "" {
		b["firstName"] = p.FirstName
	}
	if p.LastName != "" {
		b["lastName"] = p.LastName
	}
	return b
}

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
}

func (p CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(p.FirstName) < minFirstNameLen || len(p.FirstName) > maxFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name must be between %d - %d characters", minFirstNameLen, maxFirstNameLen)
	}
	if len(p.LastName) < minLastNameLen || len(p.LastName) > maxLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name must be between %d - %d characters", minLastNameLen, maxLastNameLen)
	}
	if len(p.Email) < minEmailLen || len(p.Email) > maxEmailLen || !isEmailValid(p.Email) {
		errors["email"] = "invalid email"
	}
	if len(p.Password) < minPasswordLen || len(p.Password) > maxPasswordLen {
		errors["password"] = fmt.Sprintf("password must have between %d - %d characters", minPasswordLen, maxPasswordLen)
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lastName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPassword string `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params *CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
