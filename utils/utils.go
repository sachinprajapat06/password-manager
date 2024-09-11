package utils

import (
	m "password/models"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword encrypts the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords compares a plain password with a hashed one
func ComparePasswords(hashedPassword, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}

// Function to evaluate a password based on the rules defined in the Parameter struct
func EvaluatePassword(password string) m.Parameter {
	var param m.Parameter

	// Flags to keep track of whether we found specific characteristics
	hasNum := false
	hasLower := false
	hasUpper := false
	hasSpecial := false

	// Iterate over each character in the password
	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNum = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	// Fill in the parameter struct with the validation results
	param.Have8Char = len(password) >= 8
	param.HaveNum = hasNum
	param.SmallLetter = hasLower
	param.CapitalLetter = hasUpper
	param.SpecialChar = hasSpecial

	// Evaluate "strong" password (at least 8 characters, a number, a lowercase, and an uppercase letter)
	param.Strong = param.Have8Char && param.HaveNum && param.SmallLetter && param.CapitalLetter && len(password) >= 12

	// Evaluate "super strong" password (strong password + special character)
	param.SuperStrong = param.Strong && param.SpecialChar

	return param
}
