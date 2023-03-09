package controller

import (
	"blog/config"
	"blog/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Token struct {
	Username string `json:"username"`
	JWTToken string `json:"token"`
}

var secretkey string = "secretkeyjwt"

func GenerateJWT(username string) (string, error) {
	var response Response
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		response.Message = err.Error()
		return "", err
	}
	return tokenString, nil
}
func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var response Response
	var auth Authentication
	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)
		return
	}

	var authuser model.User
	config.DB.Where("username = ?", auth.Username).First(&authuser)
	if authuser.Username == "" {
		response.Message = "Pasworrd atau email salah"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	check := CheckPasswordHash(auth.Password, authuser.Password)

	if !check {
		response.Message = "Password atau email salah"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	validToken, err := GenerateJWT(authuser.Username)
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(err)

	}

	var token Token
	token.Username = authuser.Username
	token.JWTToken = validToken
	log.Println(authuser.Username)
	log.Println(validToken)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}
func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var response Response
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.Message = err.Error()
	}
	var dbuser model.User
	config.DB.Where("username = ?", user.Username).First(&dbuser)

	//checks if email is already register or not
	if dbuser.Username != "" {
		response.Message = "Username telah digunakan "
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	user.Password, err = GeneratehashPassword(user.Password)
	if err != nil {
		log.Fatalln("error in password hash")
	}
	err = config.DB.Create(&user).Error
	if err != nil {
		response.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
	}
	response.Message = "Account Succes Created"
	response.Data = user
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)
}
func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func IsAunthenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		var response Response
		if r.Header["Token"] == nil {
			response.Message = "No Token Found"
			json.NewEncoder(w).Encode(response)
			return
		}

		var mySigningKey = []byte(secretkey)

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			response.Message = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var user model.User
			mapstructure.Decode(claims, &user)
			// username := claims["username"]
			json.NewEncoder(w).Encode(user)
			next.ServeHTTP(w, r)
		} else {
			response.Message = err.Error()
			json.NewEncoder(w).Encode(response)
			return
		}

	})
}
