package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/sessions"
)

// var Check Page

var tpl *template.Template
var Store = sessions.NewCookieStore([]byte("session"))

// template init
func init() {
	tpl = template.Must(template.ParseGlob("static/*.html"))
}

type Page struct {
	Status     bool
	Header1    interface{}
	IsLoggedin bool
	Valid      bool
}

var userData = map[string]string{
	"email":    "admin@gmail.com",
	"password": "admin123",
}

var P = Page{
	Status: false,
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	ok := Middleware(w, r)

	if ok {

		http.Redirect(w, r, "/login-submit", http.StatusSeeOther)
		return
	}
	P.Valid = Middleware(w, r)
	filename := "trial.html"
	err := tpl.ExecuteTemplate(w, filename, P)
	if err != nil {
		fmt.Println("error while parsing file", err)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if r.Method == "GET" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ok := Middleware(w, r)
		if ok {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "there is an error parsing %v", err)
		return
	}
	emails := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	if userData["email"] == emails && userData["password"] == password && r.Method == "POST" {

		session, _ := Store.Get(r, "started")

		session.Values["id"] = "Admin"
		P.Header1 = session.Values["id"]
		fmt.Println(P.Header1)
		session.Save(r, w)

		fmt.Println(session)

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		// dialog.Alert("wrong passwod")
		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return

	}

}
func Logouthandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if P.Status == true {
		session, _ := Store.Get(r, "started")
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		P.Status = false
	} else if P.Status == false {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Middleware(w http.ResponseWriter, r *http.Request) bool {
	session, _ := Store.Get(r, "started")
	if session.Values["id"] == nil {
		return false
	}
	P.Header1 = session.Values["id"]
	return true

}

func index(w http.ResponseWriter, r *http.Request) {
	ok := Middleware(w, r)
	if ok {
		P.Status = true
	}
	filenamE := "index.html"
	err := tpl.ExecuteTemplate(w, filenamE, P)
	if err != nil {
		fmt.Println("error while parsing file", err)
		return
	}

}

var Username string

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login-submit", loginHandler)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", Logouthandler)
	fmt.Println("server starts at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
