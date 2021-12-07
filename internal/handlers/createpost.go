package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type PostInfo struct {
	Title string
	Text  string
}

var PostData PostInfo

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		PostData.Title = r.FormValue("title")
		PostData.Text = r.FormValue("text")

		http.Redirect(w, r, "/", 302)
	} else {

		tpl, err := template.ParseFiles("./templates/createpost.html", "./templates/base.html")
		checkErr(err)
		tpl.ExecuteTemplate(w, "createpost.html", nil)

	}
}
