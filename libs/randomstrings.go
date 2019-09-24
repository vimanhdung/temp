// thanks to @elithrar for the code to create the secret token!
// source: https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/
package libs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

// GenerateRandomBytes returns securely generated random bytes. 
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GenerateMd5String returns a Md5 string encoded
// It will return an error if the system's secure random
func GenerateMd5String() (string) {
	md5Str, _ := GenerateRandomString(32)
	hasher := md5.New()
	hasher.Write([]byte(md5Str))
	return hex.EncodeToString(hasher.Sum(nil))
}