package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

type SearchData struct {
	search     bool
	searchText string
}

func ProductsHandle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	products := getProductsByName("")
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

func SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	name := r.FormValue("productName")
	products := getProductsByName(name)

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
func getProductsByName(name string) []Product {
	db, err := sql.Open("mysql", "root:password@/world")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	result, err := db.Query("SELECT * FROM products WHERE car_name LIKE ?;", "%"+name+"%")
	if err != nil {
		log.Println(err)
	}

	products := []Product{}
	for result.Next() {
		var p Product
		err = result.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price)
		if err != nil {
			log.Println(err)
		}
		products = append(products, p)
	}

	return products
}
func filtredProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	r.ParseForm()

	minPriceStr := r.FormValue("minPrice")
	minPrice, err := strconv.Atoi(minPriceStr)
	if err != nil {
		log.Println(err)
	}

	maxPriceStr := r.FormValue("maxPrice")
	maxPrice, err := strconv.Atoi(maxPriceStr)
	if err != nil {
		log.Println(err)
	}

	products, err := getFilteredProducts(database, minPrice, maxPrice)
	if err != nil {
		log.Println(err)
	}

	toInData := toIndexData{
		Products: products,
		Username: savUsername,
	}
	tpl.ExecuteTemplate(w, "index.html", toInData)
}

func getFilteredProducts(db *sql.DB, minPrice, maxPrice int) ([]Product, error) {
	// Prepare the SQL query
	query := "SELECT id, car_name,details, price FROM products WHERE price >= ? AND price <= ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the query
	rows, err := stmt.Query(minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the results
	products := []Product{}

	// Iterate over the rows
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

var tpl *template.Template

func main() {

	tpl, _ = template.ParseGlob("templates/*.html")

	db, err := sql.Open("mysql", "root:password@/world")
	if err != nil {
		log.Println(err)
	}
	defer db.Close()

	database = db

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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
	}

	for _, route := range routes {
		http.HandleFunc(route.path, route.handler)
	}

	products, err := getFilteredProducts(db, 50, 100)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(products)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
