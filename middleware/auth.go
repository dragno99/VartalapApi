package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mySigningKey = []byte("VartalapChatApi3838jbdbdj993b3idw9dw2e")

func GenerateJWT(userId primitive.ObjectID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 60000).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				json.NewEncoder(w).Encode(bson.M{
					"message": err.Error(),
				})
				return
			}
			if token.Valid {
				claims := token.Claims.(jwt.MapClaims)
				userId := claims["userId"].(string)
				r.Header.Add("userId", userId)
				endpoint(w, r)
			} else {
				json.NewEncoder(w).Encode(bson.M{
					"message": "Token is not valid",
				})
			}
		} else {
			json.NewEncoder(w).Encode(bson.M{
				"message": "Not authorized",
			})
		}
	})
}
