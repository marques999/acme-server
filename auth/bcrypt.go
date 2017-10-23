package auth

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}

func GeneratePassword(rawPassword string) (string, error) {

	if password, ex := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost); ex == nil {
		return string(password), nil
	} else {
		return "", ex
	}
}
