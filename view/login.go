package view

import (
	"fmt"
	"html/template"
	"midterm1/config"
	"midterm1/controller"
	"midterm1/models"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, tpl *template.Template) {
	db, err := config.LoadDB()
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", "")
		return
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var pass string
	var userid int64
	stmt := "SELECT userid, password FROM users WHERE name = ?"
	row := db.QueryRow(stmt, username)
	err = row.Scan(&userid, &pass)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	if password == pass {
		products := controller.GetProductsByName("")

		data := models.ToIndexData{
			Username: username,
			UserId:   userid,
			Products: products}
		tpl.ExecuteTemplate(w, "index.html", data)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
	fmt.Println("incorrect password")

}
