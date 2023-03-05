package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	Id       int
	fullname string
	password string
}

var database *sql.DB

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
		http.ServeFile(w, r, "templates/create.html")
	}
}

func main() {

	db, err := sql.Open("mysql", "root:password@/testdb")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	if err != nil {
		log.Println(err)
	}
	database = db
	http.HandleFunc("/create", CreateHandler)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}
