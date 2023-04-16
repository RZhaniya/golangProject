package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/go-humble/locstor"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gorilla/mux"
)

type toIndexData struct {
	Username string
	Products []Product
	UserId   int64
}
type Product struct {
	Id       int
	Car_name string
	Details  string
	Price    int
	Rate     string
}

var database *sql.DB
var toInData toIndexData
var savUsername string
var userid int64
var tmpl *template.Template

type SearchData struct {
	search     bool
	searchText string
}
type Comment struct {
	Name    string
	Comment string
}
type ProductPage struct {
	Product  Product
	Comments []Comment
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
		UserId:   userid,
	}

	// tmpl = template.Must(template.ParseFiles("./templates/*"))
	nerr := tpl.Execute(w, toInData)

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
		UserId:   userid,
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
	// fmt.Println("password:", password, "\npswdLength:", len(password))

	stmt := "SELECT userid FROM users WHERE name = ?"
	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
	row := db.QueryRow(stmt, username)

	err = row.Scan(&userid)

	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "registration.html", "username already taken")
		return
	}
	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO users (name, password) VALUES (?, ?);")

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
	userid = lastIns
	savUsername = username
	products := getProductsByName("")

	data := toIndexData{
		Username: savUsername,
		UserId:   userid,
		Products: products}
	tpl.ExecuteTemplate(w, "index.html", data)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", "")
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var pass string
	stmt := "SELECT userid, password FROM users WHERE name = ?"
	row := db.QueryRow(stmt, username)
	err = row.Scan(&userid, &pass)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	if password == pass {
		products := getProductsByName("")

		data := toIndexData{
			Username: username,
			UserId:   userid,
			Products: products}
		tpl.ExecuteTemplate(w, "index.html", data)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("incorrect password")

}
func getProductsByName(name string) []Product {
	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
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
		err = result.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate)
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
	query := "SELECT id, car_name,details, price,rate FROM products WHERE price >= ? AND price <= ?"
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
		if err := rows.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
func getProduct(productId string) (Product, error) {
	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
	if err != nil {
		return Product{}, err
	}
	defer db.Close()

	result := db.QueryRow("SELECT * FROM products WHERE id = ?", productId)
	var p Product
	err = result.Scan(&p.Id, &p.Car_name, &p.Details, &p.Price, &p.Rate)
	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func getComments(productId string) ([]Comment, error) {
	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
	if err != nil {
		return []Comment{}, err
	}
	defer db.Close()

	res, err := db.Query("SELECT u.name, c.comment FROM comments c join users u on u.userid=c.userid WHERE productId = ?", productId)
	if err != nil {
		return []Comment{}, err
	}

	comments := []Comment{}
	for res.Next() {
		var c Comment
		err = res.Scan(&c.Name, &c.Comment)
		if err != nil {
			return []Comment{}, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func renderProductPage(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productId := params["id"]

	p, err := getProduct(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
		return
	}

	comments, err := getComments(productId)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	data := ProductPage{
		Product:  p,
		Comments: comments,
	}

	tpl.ExecuteTemplate(w, "product.html", data)
}

func sendComment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside")
	var result string
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "product.html", nil)
		return
	}

	commentText := r.FormValue("commentText")
	fmt.Println(commentText)
	if len(commentText) == 0 {
		result = "write some comment"
		tpl.ExecuteTemplate(w, "product.html", result)
		return
	} else {

		userId := r.FormValue("userId")
		productId := r.FormValue("productId")
		var insertStmt *sql.Stmt
		insertStmt, err2 := database.Prepare("INSERT INTO comments (productid,userid, comment) VALUES (?, ?,?);")
		// fmt.Println(userId)
		if err2 != nil {
			fmt.Println("error preparing statement:", err2)
			tpl.ExecuteTemplate(w, "index.html", "there was a problem registering account")
			return
		}
		defer insertStmt.Close()
		var result sql.Result

		result, err2 = insertStmt.Exec(productId, userId, commentText)
		lastIns, _ := result.LastInsertId()
		fmt.Println("lastIns comment:", lastIns)
		if err2 != nil {
			fmt.Println("error inserting new user")
			tpl.ExecuteTemplate(w, "registration.html", "there was a problem registering account")
			return
		}
		p, err := getProduct(productId)
		fmt.Println(productId)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve product", http.StatusInternalServerError)
			return
		}

		comments, err := getComments(productId)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
			return
		}

		data := ProductPage{
			Product:  p,
			Comments: comments,
		}

		http.Redirect(w, r, "/product:"+productId, http.StatusSeeOther)
		tpl.ExecuteTemplate(w, "product.html", data)
	}
}
func sendRating(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() // Parses the request body
	rating := r.Form.Get("rating")
	productId := r.Form.Get("productId")
	userId := r.Form.Get("userId") // x will be "" if parameter is not set
	fmt.Println(rating + " " + productId + " " + userId)

	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")

	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO ratings (rate, productId,userId) VALUES (?, ?,?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		return
	}
	defer insertStmt.Close()
	var result sql.Result
	result, err = insertStmt.Exec(rating, productId, userid)
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

	db, err := sql.Open("mysql", "root:@(localhost:3306)/world")
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
	// r.HandleFunc("/books/{id}", getBook).Methods("GET")
	// r.HandleFunc("/books", createBook).Methods("POST")
	// r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	// r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	for _, route := range routes {
		r.HandleFunc(route.path, route.handler)
	}

	log.Fatal(http.ListenAndServe(":8080", r))

	fmt.Println("Server is listening...")
	// http.ListenAndServe(":8080", nil)
}
