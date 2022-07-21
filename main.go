package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id                    uint16
	Title, Anons, FulText string
}

var posts = []Article{}
var showPost = Article{}

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

func HandleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	//указывем что обработка всх адресов будет через горилла роутер
	http.Handle("/", rtr)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	http.ListenAndServe(":8080", nil)

}

func main() {
	HandleFunc()

}
