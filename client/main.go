package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

//vqr mySigningKey = os.Get("MY_JWT_TOKEN")
var mySigningKey = []byte("quakegodmode")

//better to pick this up from env variables
// set MY_JWT_TOKEN=quakegodmode (has to be set in cli)

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "tosha"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("something went wrong %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

func index(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/", nil)
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

func create(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/create", nil)
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "Errors: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println("Client create func")
	fmt.Fprintf(w, string(body))

}

func save_article(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/save_article", nil)
	req.Header.Set("Token", validToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "Errors: %s", err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	fmt.Println("Client save article func")
	fmt.Fprintf(w, string(body))

}

func loginPage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/loginn", nil)
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

func check_password(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/check_password", nil)
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

func show_post(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "http://localhost:8081/post/{id:[0-9]+}", nil)
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
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)
	http.HandleFunc("/loginn", loginPage)
	http.HandleFunc("/check_password", check_password)
	http.HandleFunc("/post/{id:[0-9]+}", show_post)

}

func main() {
	fmt.Println("My Client 8080")

	HandleRequests()

	http.ListenAndServe(":8080", nil)

}
