package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthHelper struct{}

// hash password
func (h *AuthHelper) HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedPassword)
}

// verify password
func (h *AuthHelper) VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "email or password is incorrect"
		check = false
	}

	return check, msg
}
