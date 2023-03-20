package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type toIndexData struct {
	Username string
	Products []Product
}
type Product struct {
	Id       int
	Car_name string
	Details  string
	Price    int
}

var database *sql.DB
var toInData toIndexData
var savUsername string

func ProductsHandle(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	result, err := database.Query("select * from products")
	if err != nil {
		log.Println(err)
	}

	products := []Product{}

	for result.Next() {
		var P Product
		err = result.Scan(&P.Id, &P.Car_name, &P.Details, &P.Price)
		products = append(products, P)
		if err != nil {
			panic(err)
		}

	}

	toInData := toIndexData{
		Products: products,
		Username: savUsername,
	}

	var tmpl = template.Must(template.ParseFiles("./templates/index.html"))
	nerr := tmpl.Execute(w, toInData)

	if nerr != nil {
		log.Println(nerr)
	}

}

type SearchData struct {
	search     bool
	searchText string
}

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	name := "%" + r.FormValue("productName") + "%"
	fmt.Println("name:", name)
	db, err := sql.Open("mysql", "root:password@/world")

	result, err := db.Query("SELECT * FROM products WHERE car_name like ?;", name)
	products := []Product{}
	for result.Next() {
		var P Product
		err = result.Scan(&P.Id, &P.Car_name, &P.Details, &P.Price)
		products = append(products, P)
		if err != nil {
			panic(err)
		}

	}
	toInData := toIndexData{
		Products: products,
		Username: savUsername,
	}
	tpl.ExecuteTemplate(w, "index.html", toInData)

}

// registerHandler serves form for registring new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "registration.html", "")
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("password:", password, "\npswdLength:", len(password))

	stmt := "SELECT id FROM users WHERE username = ?"
	db, err := sql.Open("mysql", "root:password@/world")
	row := db.QueryRow(stmt, username)
	var uID int
	err = row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "registration.html", "username already taken")
		return
	}
	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (username, password) VALUES (?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()
	var result sql.Result

	result, err = insertStmt.Exec(username, password)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
		return
	}
	savUsername = username
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:password@/world")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", "")
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var pass string
	stmt := "SELECT password FROM users WHERE username = ?"
	row := db.QueryRow(stmt, username)
	err = row.Scan(&pass)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	if password == pass {
		savUsername = username
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("incorrect password")

}

var tpl *template.Template

func main() {
	tpl, _ = template.ParseGlob("templates/*.html")

	db, err := sql.Open("mysql", "root:password@/world")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err != nil {
		log.Println(err)
	}
	database = db
	http.HandleFunc("/", ProductsHandle)
	http.HandleFunc("/search", SearchPage)
	http.HandleFunc("/products", ProductsHandle)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
