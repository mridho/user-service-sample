package authentication

import (
	"crypto/rsa"
	"time"

	"user-service-sample/config"
	"user-service-sample/repository"

	"github.com/golang-jwt/jwt/v5"
)

var (
	tokenTTL = 72 * time.Hour
)

type jwtCustomClaims struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
	jwt.RegisteredClaims
}

func GenerateSignedToken(cfg config.SecretConfig, user repository.User) (token string, err error) {
	// Set custom claims
	claims := &jwtCustomClaims{
		Id:       user.Id,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
		},
	}

	// Create JWT with claims
	jwtWithClaim := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	privateKey, err := readRSAPrivateKey([]byte(cfg.RsaPrivatePem))
	if err != nil {
		return "", err
	}

	// Generate encoded token and send it as response.
	return jwtWithClaim.SignedString(privateKey)
}

func readRSAPrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {
	// privateKeyBytes, err := os.ReadFile(privateKeyPath)
	// if err != nil {
	// 	return nil, err
	// }

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
