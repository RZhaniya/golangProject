package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"midterm1/config"
	"midterm1/controller"
	"midterm1/models"
	"midterm1/view"
	"net/http"

	"github.com/go-humble/locstor"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

var database *sql.DB
var toInData models.ToIndexData
var savUsername string
var userid int64
var tmpl *template.Template

type Comment struct {
	Name    string
	Comment string
}

func ProductsHandle(w http.ResponseWriter, r *http.Request) {
	view.ProductsHandle(w, r, tpl)

}
func SearchPage(w http.ResponseWriter, r *http.Request) {
	controller.SearchPage(w, r, tpl)
}

// registerHandler serves form for registring new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	view.RegisterHandler(w, r, tpl)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	view.LoginHandler(w, r, tpl)

}
func filtredProduct(w http.ResponseWriter, r *http.Request) {
	controller.FiltredProduct(w, r, tpl)
}

func renderProductPage(w http.ResponseWriter, r *http.Request) {
	controller.RenderProductPage(w, r, tpl)
}

func sendComment(w http.ResponseWriter, r *http.Request) {
	controller.SendComment(w, r, tpl)
}
func sendRating(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // Parses the request body
	rating := r.Form.Get("rating")
	productId := r.Form.Get("productId")
	userId := r.Form.Get("userId") // x will be "" if parameter is not set

	db, err := config.LoadDB()

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO ratings (rate, productId,userId) VALUES (?, ?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	result, err = insertStmt.Exec(rating, productId, userId)
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	lastIns, _ := result.LastInsertId()
	fmt.Print(lastIns)
}

var tpl *template.Template

func main() {
	if err := locstor.SetItem("userId", "1"); err != nil {
		fmt.Println(err)
	}
	tpl, _ = template.ParseGlob("templates/*.html")

	db, err := config.LoadDB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	database = db

	routes := []struct {
		path    string
		handler http.HandlerFunc
	}{
		{path: "/", handler: ProductsHandle},
		{path: "/search", handler: SearchPage},
		{path: "/products", handler: ProductsHandle},
		{path: "/register", handler: registerHandler},
		{path: "/login", handler: loginHandler},
		{path: "/filtred", handler: filtredProduct},
		{path: "/product:{id}", handler: renderProductPage},
		{path: "/sendComment", handler: sendComment},
		{path: "/ratings", handler: sendRating},
	}
	r := mux.NewRouter()
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	for _, route := range routes {
		r.HandleFunc(route.path, route.handler)
	}

	log.Fatal(http.ListenAndServe(":8080", r))

	fmt.Println("Server is listening...")
	// http.ListenAndServe(":8080", nil)
}
