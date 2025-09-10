package helpers

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateVerificationCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)
	return code, nil
}