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
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("⛔️ Something Went Wrong: %s", err.Error())
	}
	return tokenString, nil
}

func readJWT(unparsedToken string) (*h.User, error) {
	godotenv.Load()
	var mySigningKey = []byte(os.Getenv("SIGN_KEY"))
	token, err := jwt.Parse(unparsedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("⛔️ There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("⛔️ Your Token has been expired")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user h.User
		user.Email = fmt.Sprintf("%v", claims["email"])
		return &user, nil
	}
	return nil, fmt.Errorf("⛔️ Something went wrong during token read")
}