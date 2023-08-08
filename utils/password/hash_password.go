package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const (
	saltByteLength = 10
)

func SaltAndHashPassword(password string) (passwordHash, salt string) {
	// generate salt
	salt = generateRandomString(saltByteLength)

	// generate hash
	passwordHash = createHash(password, salt)

	return passwordHash, salt
}

func CheckPassword(password, passwordHash, salt string) bool {
	return passwordHash == createHash(password, salt)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func createHash(password, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return hex.EncodeToString(hasher.Sum(nil))
}
