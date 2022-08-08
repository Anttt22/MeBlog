package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User{
	Username: "1",
	Password: "1",
}

type Article struct {
	Id                    uint16
	Title, Anons, FulText string
}

var posts = []Article{}
var showPost = Article{}
var mySigningKey = []byte("quakegodmode")

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
		//fmt.Println(fmt.Sprintf("Id: %d, title: %s , anons: %s , text: %s", post.id, post.Title, post.Anons, post.FulText))
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
		http.Redirect(w, r, "/", 301)
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
		//fmt.Println(fmt.Sprintf("Id: %d, title: %s , anons: %s , text: %s", post.id, post.Title, post.Anons, post.FulText))
		t.ExecuteTemplate(w, "show", showPost)
	}

	// id := vars["id"]
	// response := fmt.Sprintf("id= %s", id)
	// fmt.Fprint(w, response)
}

func loginPage(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprint(w, "Loginn page")

	t, err := template.ParseFiles("templates/login.html", "templates/header.html", "templates/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "login", nil)

}

func check_password(w http.ResponseWriter, r *http.Request) {
	//to be replaced by db later - hardcodeedfor now

	var dbname string = "1"
	var dbpass string = "1"

	name_to_check := r.FormValue("userlogin")
	pass_to_check := r.FormValue("userpassword")

	// if name_to_check == "" || pass_to_check == "" {
	// 	fmt.Fprintf(w, "not all data in fields pass or login filled")
	// }

	if dbname == name_to_check && dbpass == pass_to_check {
		http.Redirect(w, r, "/", 301)
	} else {
		fmt.Fprintf(w, "not all data in fields pass or login filled")
	}

	//http.Redirect(w, r, "/create", 301)

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

func HandleFunc() {
	//port := os.Getenv("PORT")
	//log.Print("Listen on :" + port)

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/loginn", loginPage)
	rtr.HandleFunc("/check_password", check_password).Methods("POST")
	rtr.Handle("/create", checkAuth(create)).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	//указывем что обработка всх адресов будет через горилла роутер
	http.Handle("/", rtr)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8081", nil)
	//log.Fatal(http.ListenAndServe(":"+port, nil))

}

func main() {
	HandleFunc()

}
