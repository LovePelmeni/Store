package authentication

import (
	"errors"
	"log"

	"os"
	"time"

	"github.com/LovePelmeni/Store/models"
	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secretKey = os.Getenv("JWT_AUTH_SECRET_KEY")
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	WarnLogger  *log.Logger
)

// Jwt Exceptions..
func InvalidJwt() error {
	return errors.New("Invalid JWT Token")
}

func InvalidJwtKey() error {
	return errors.New("Invalid JWT Secret Key")
}

func InvalidJwtSignature() error {
	return errors.New("Invalid JWT Signature.")
}

func JwtDecodeError() error {
	return errors.New("Failed to Decode JWT Token.")
}

type JwtToken struct {
	jwt.StandardClaims

	Username string
	Email    string
}

func CreateJwtToken(Username string, Email string) (string, error) {

	ExpirationTime := time.Now().Add(10000 * time.Minute)
	claims := JwtToken{
		Username: Username,
		Email:    Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: ExpirationTime.Unix(),
		},
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, error := newToken.SignedString(secretKey)
	if error != nil {
		ErrorLogger.Println("Failed to Stringify JWT Token.")
		return "", error
	}
	return stringToken, nil
}

type JwtValidator struct {
	Token string
}

type DecodedJwtData struct {
	Username string
	Email    string
}

func CheckValidJwtToken(token string) error {

	var customer models.Customer
	DecodedData := &JwtToken{}
	_, error := jwt.ParseWithClaims(token, DecodedData,
		func(token *jwt.Token) (interface{}, error) { return secretKey, nil })

	if error != nil {
		InfoLogger.Println("Invalid Jwt Token")
		return InvalidJwt()
	}

	if customer := models.Database.Table("customers").Where("Username = ? AND Email = ?",
		DecodedData.Username,
		DecodedData.Email).First(&customer); customer.Error != nil {
		return InvalidJwt()
	}
	return nil
}

func GetCustomerJwtCredentials(token string) (map[string]string, error) {

	DecodedData := &JwtToken{}
	_, error := jwt.ParseWithClaims(token, DecodedData,
		func(token *jwt.Token) (interface{}, error) { return secretKey, nil })

	if error != nil {
		InfoLogger.Println("Invalid Jwt Token")
		return nil, InvalidJwt()
	}
	return map[string]string{"username": DecodedData.Username,
		"email": DecodedData.Email}, nil
}
