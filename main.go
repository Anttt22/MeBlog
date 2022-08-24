package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Article struct {
	Id                    uint16
	Title, Anons, FulText string
}

var posts = []Article{}
var showPost = Article{}
var mySigningKey = []byte("quakegodmode")

func checkAuthviaCookie(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		c, err := r.Cookie("token")

		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// For any other type of error, return a bad request status
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tknStr := c.Value

		w.Header().Set("Connection", "close")
		defer r.Body.Close()

		if c.Value != "" {
			token, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("errror in header token parsing")
				}
				return mySigningKey, nil
			})

			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				fmt.Println(err.Error())
				return

				//!!! problrem here made it redirect to login page...
			}

			if token.Valid {
				endpoint(w, r)
			} else {
				fmt.Println("token not valid")
			}

		} else {

			fmt.Fprintf(w, "Not Authorizeddd from cookie")

		}
	})
}

func checkAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Connection", "close")
		defer r.Body.Close()
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("errror in header token parsing")
				}
				return mySigningKey, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Header().Add("Content-Type", "application/json")
				return
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorizeddd")

		}
	})
}

func index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:Super_22@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FulText)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)

	}
	//не просто execute потому что подключим динаическое отображние шаблонов
	t.ExecuteTemplate(w, "index", posts)
	//index- name of block. - {{index}}
}

func create(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	//не просто execute потому что подключим динаическое отображние шаблонов
	t.ExecuteTemplate(w, "create", nil)
	//index- name of block. - {{index}}
}

func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "not all data filled")
	} else {

		db, err := sql.Open("mysql", "root:Super_22@tcp(127.0.0.1:3306)/golangdb")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES ('%s', '%s', '%s')", title, anons, full_text))
		if err != nil {
			panic(err)
		}
		defer insert.Close()

		http.Redirect(w, r, "/create", 301)
	}
}

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	t, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:Super_22@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}
	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FulText)
		if err != nil {
			panic(err)
		}
		showPost = post
		t.ExecuteTemplate(w, "show", showPost)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "login", nil)
}

func check_hashed_password(w http.ResponseWriter, r *http.Request) {
	var hashed_pass_from_db string

	name_to_check := r.FormValue("userlogin")
	pass_to_check := r.FormValue("userpassword")
	db, err := sql.Open("mysql", "root:Super_22@tcp(127.0.0.1:3306)/golangdb")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT `password` FROM `users` WHERE `name` = '%s'", name_to_check))
	if err != nil {
		panic(err)
	}

	for res.Next() {
		err = res.Scan(&hashed_pass_from_db)
		if err != nil {
			panic(err)
		}
	}
	if bcrypt.CompareHashAndPassword([]byte(hashed_pass_from_db), []byte(pass_to_check)) == nil {
		validToken, err := GenerateJWT()
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: validToken,
		})
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

	} else {
		fmt.Fprintf(w, "!!relogin with a right credentials!!")
	}
}

//function hashing password (pass) to be used with a register page later
func to_hash_password(pass string) []byte {
	password := []byte(pass)
	cost := 10
	hash, _ := bcrypt.GenerateFromPassword(password, cost)
	return hash
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "tosha"
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("fail with JWT token generation %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func HandleFunc() {
	//port := os.Getenv("PORT")  //for remote server
	//log.Print("Listen on :" + port)  //for remote server

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/loginn", loginPage)
	rtr.HandleFunc("/check_password", check_hashed_password)
	rtr.HandleFunc("/save_article", save_article)
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")
	rtr.Handle("/create", checkAuthviaCookie(create)).Methods("GET")

	//routing via Gorilla router
	http.Handle("/", rtr)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8081", nil)

	//log.Fatal(http.ListenAndServe(":"+port, nil))   //for remote server
}

func main() {
	HandleFunc()

}

//sozdat bd users s id imenem i hasged password

//esli token ne valid to create post  - no page found> inly after login is going to work
//dobavit avtorizaciu with login and pass i hraneniev bd
// sdelat otdelnij server dlya proverki parolya i logina
// bd users and passwords // mozgno ispolzovat tu zge samuyu bd tolko otdelnyu tablicy
