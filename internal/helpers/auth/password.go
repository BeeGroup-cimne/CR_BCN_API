package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Generate a hashed version of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	// Compare the hashed password with the plain text password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // If err == nil, it means the passwords match
}
