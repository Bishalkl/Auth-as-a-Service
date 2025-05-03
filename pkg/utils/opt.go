package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateOTP generates a secure random numeric OTP of the given length
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	otp := ""

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp += string(digits[num.Int64()])
	}

	return otp, nil
}
