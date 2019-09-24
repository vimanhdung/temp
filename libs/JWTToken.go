package libs

import (
	"crypto/rsa"
	"errors"
	"indetail/conf"
	"indetail/models"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// EasyToken is an Struct to encapsulate username and expires as parameter
type EasyToken struct {
	// Username is the name of the user
	Username string
	// JTI of token
	Jti string
}

// https://gist.github.com/cryptix/45c33ecf0ae54828e63b
// location of the files used for signing and verification
const (
	privKeyPath = "libs/keys/rsakey.pem"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "libs/keys/rsakey.pem.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	verifyKey    *rsa.PublicKey
	mySigningKey *rsa.PrivateKey
)

func init() {
	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Println(err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Println(err)
	}

	signBytes, err := ioutil.ReadFile(privKeyPath)

	if err != nil {
		log.Println(err)
	}

	mySigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Println(err)
	}
}

// GetToken is a function that exposes the method to get a simple token for jwt
func (e EasyToken) GetToken(expires time.Duration) (string, error) {
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(expires).Unix(), //time.Unix(c.ExpiresAt, 0)
		Issuer:    e.Username,
		Id:        e.Jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
	}

	return tokenString, err
}

// RefreshToken is s function that exposes the method to get a refresh token for jwt
func (e EasyToken) RefreshToken(tokenString string, user *models.AdminAccounts) (string, error) {
	// get jti for token
	jti := GenerateMd5String()
	oldToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldClaims, _ := oldToken.Claims.(jwt.MapClaims)

	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(conf.TokenExpires).Unix(), //time.Unix(c.ExpiresAt, 0)
		Issuer:    oldClaims["iss"].(string),
		Id:        jti,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshToken, err := token.SignedString(mySigningKey)
	if err != nil {
		log.Println(err)
	}

	// Destroy token
	// Save new token in DB
	if b := models.DestroyTokenAndSaveNewToken(oldClaims["jti"].(string), jti, user); !b {
		log.Println(b)
	}

	return refreshToken, err
}

// ValidateToken get token strings and return if is valid or not
func (e EasyToken) ValidateToken(tokenString string) (bool, string, string, error) {
	// Token from another example.  This token is expired
	//var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJleHAiOjE1MDAwLCJpc3MiOiJ0ZXN0In0.HE7fK0xOQwFEr4WDgRWj4teRPZ6i3GLwD5YCm6Pwu_c"
	if tokenString == "" {
		return false, "", "", errors.New("token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	if token == nil {
		return false, "", "", errors.New("not work")
	}

	if token.Valid {
		//"You look nice today"
		claims, _ := token.Claims.(jwt.MapClaims)
		//var user string = claims["username"].(string)
		iss := claims["iss"].(string)
		jti := claims["jti"].(string)
		return true, iss, jti, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return false, "", "", errors.New("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return false, "", "", errors.New("Timing is everything")
		} else {
			//"Couldn't handle this token:"
			return false, "", "", err
		}
	} else {
		//"Couldn't handle this token:"
		return false, "", "", err
	}
}
