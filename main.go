package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type session struct {
	username  string
	userLevel int
	isLogged  bool
	token     string
	email     string
}

type mainPage struct {
	Title string
}

type page struct {
	_SESSION session
	MainPage mainPage
}

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {

	http.HandleFunc("/", index)
	http.HandleFunc("/process", process)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	if checkURL(req.URL.Path) {
		/*
			username  string
			userLevel int
			isLogged  bool
			token     string
			email     string
		*/
		sessionData := session{"guest", 0, false, "0", "guest@gmail.com"}

		mainPageData := mainPage{"Forum"}

		PageData := page{sessionData, mainPageData}

		fmt.Println(PageData.MainPage.Title)
		tpl.ExecuteTemplate(os.Stdout, "index.gohtml", PageData)
		err := tpl.ExecuteTemplate(w, "index.gohtml", PageData)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("404 page not found:", req.URL.Path)
		err := tpl.ExecuteTemplate(w, "404.gohtml", nil)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}

func checkURL(urlGiven string) bool {
	res := true
	check := ""
	if strings.Contains(urlGiven, ".gohtml") {
		check = strings.Split(urlGiven, ".gohtml")[0]
	} else {
		check = urlGiven
	}

	switch check {
	case "/index":
	case "/":
	case "/process":
	default:
		if _, err := os.Stat(urlGiven); os.IsNotExist(err) {
			res = false
		}
	}

	return res
}

func process(w http.ResponseWriter, req *http.Request) {

	sendTo := "index.gohtml"

	if req.Method == http.MethodPost {
		subVal := req.FormValue("subVal")

		switch subVal {
		case "subAsciiArt":
			//subAsciiArtProc(w, req)
			break
		default:
			log.Println(subVal)

			err := tpl.ExecuteTemplate(w, sendTo, nil)

			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
	}
}
