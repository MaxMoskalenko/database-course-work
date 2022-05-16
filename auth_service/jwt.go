package auth_service

import (
	h "database-course-work/helpers"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func generateJWT(user *h.User) (string, error) {
	godotenv.Load()
	var mySigningKey = []byte(os.Getenv("SIGN_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("something Went Wrong: %s", err.Error())
	}
	return tokenString, nil
}

func generateCompanyJWT(company *h.Company) (string, error) {
	godotenv.Load()
	var mySigningKey = []byte(os.Getenv("SIGN_KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["tag"] = company.Tag
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("something Went Wrong: %s", err.Error())
	}
	return tokenString, nil
}

func ReadJWT(unparsedToken string) (*h.User, error) {
	godotenv.Load()
	var mySigningKey = []byte(os.Getenv("SIGN_KEY"))
	token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user h.User
		user.Email = fmt.Sprintf("%v", claims["email"])
		return &user, nil
	}
	return nil, fmt.Errorf("something went wrong during token read")
}

func ReadCompanyJWT(unparsedToken string) (*h.Company, error) {
	godotenv.Load()
	var mySigningKey = []byte(os.Getenv("SIGN_KEY"))
	token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var company h.Company
		company.Tag = fmt.Sprintf("%v", claims["tag"])
		return &company, nil
	}
	return nil, fmt.Errorf("something went wrong during token read")
}
