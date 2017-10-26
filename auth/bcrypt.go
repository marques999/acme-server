package auth

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(hashedPassword string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}

func KamikazePassword(password string) string {
	hashedPassword, _ := GeneratePassword(password)
	return hashedPassword
}

func GeneratePassword(password string) (string, error) {

	if password, ex := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); ex == nil {
		return string(password), nil
	} else {
		return "", ex
	}
}