package authentication

import (
	"crypto/rsa"
	"errors"
	"strings"

	"user-service-sample/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
)

const (
	AuthHeaderKey = "Authorization"
)

var (
	ErrInvalidAutheticatonHeader = errors.New("invalid Authentication header")
	ErrAccessingClaims           = errors.New("error accessing claims")
	ErrDecodingClaims            = errors.New("error decoding claims")
)

func VerifyToken(ctx echo.Context, cfg config.SecretConfig) (cc *jwtCustomClaims, err error) {
	authHeaderVal := strings.TrimSpace(ctx.Request().Header.Get(AuthHeaderKey))
	authHeaderSplit := strings.Split(authHeaderVal, " ")
	if len(authHeaderSplit) != 2 {
		return nil, ErrInvalidAutheticatonHeader
	}
	receivedToken := authHeaderSplit[1]

	publicKey, err := readRSAPublicKey([]byte(cfg.RsaPublicPem))
	if err != nil {
		return nil, err
	}

	// Parse and validate the token
	token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidAutheticatonHeader
	}

	// Access the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrAccessingClaims
	}

	// Safely parse the claims into custom claims struct
	cc = new(jwtCustomClaims)
	if err := mapstructure.Decode(claims, cc); err != nil {
		return nil, ErrDecodingClaims
	}

	return cc, nil
}

func readRSAPublicKey(publicKeyBytes []byte) (*rsa.PublicKey, error) {
	// publicKeyBytes, err := os.ReadFile(publicKeyPath)
	// if err != nil {
	// 	return nil, err
	// }

	privateKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
