package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	Id       int
	Car_name string
	Details  string
	Price    int
}

var database *sql.DB

func ProductsHandle(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	result, err := database.Query("select * from products")
	if err != nil {
		log.Println(err)
	}
	product := Product{}
	products := []Product{}

	for result.Next() {
		var id int
		var car_name string
		var details string
		var price int

		err = result.Scan(&id, &car_name, &details, &price)

		product.Id = id
		product.Car_name = car_name
		product.Details = details
		product.Price = price
		products = append(products, product)

		if err != nil {
			panic(err)
		}

	}

	var tmpl = template.Must(template.ParseFiles("./templates/products.html"))
	nerr := tmpl.Execute(w, products)

	if nerr != nil {
		log.Println(nerr)
	}
	// http.ServeFile(w, r, "./templates/registration.html")

}
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		model := r.FormValue("fullname")
		company := r.FormValue("password")

		_, err = database.Exec("insert into testdb.usertest (id,fullname,password) values (?, ?, ?)",
			1, model, company)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "./templates/craete.html")
	}

}
func pageHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("working page")

}
func main() {

	db, err := sql.Open("mysql", "root:password@/world")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err != nil {
		log.Println(err)
	}
	database = db
	http.HandleFunc("/", ProductsHandle)
	http.HandleFunc("/create", CreateHandler)
	http.HandleFunc("/products", ProductsHandle)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
