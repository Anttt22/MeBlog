package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

//better to pick this up from env variables
// set MY_JWT_TOKEN=quakegodmode (has to be set in cli)
//vqr mySigningKey = os.Get("MY_JWT_TOKEN")
var mySigningKey = []byte("quakegodmode")

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User{
	Username: "1",
	Password: "1",
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/create", nil)
	req.Header.Set("Token", checkLogin(u))
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "Errors: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, string(body))

}

func checkLogin(u User) string {
	//	fmt.Println("\ncp checl login 1")
	if user.Username != u.Username || user.Password != u.Password {
		fmt.Println("user not found")
		err := "error"
		return err
	}
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Println(err)
	}
	return validToken
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "tosha"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("something wrong in JWT token generation %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func homePage(w http.ResponseWriter, r *http.Request) {

	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/stocks", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "Errors: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprintf(w, string(body))

}

func HandleRequests() {
	//http.HandleFunc("/", homePage)
	http.HandleFunc("/login", loginPage)

}

func main() {
	fmt.Println("My Client")

	HandleRequests()

	http.ListenAndServe(":8080", nil)

}
