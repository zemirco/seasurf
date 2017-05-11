package main

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type user struct {
	Name string
	Age  int
}

const (
	sessionName         = "session-name"
	sessionValueKeyUser = "user"
)

var (
	loginTemplate    *template.Template
	settingsTemplate *template.Template
	store            = sessions.NewCookieStore([]byte("something-very-secret"))
)

func init() {
	gob.Register(&user{})
	loginTemplate = template.Must(template.ParseFiles("base.html", "login.html"))
	settingsTemplate = template.Must(template.ParseFiles("base.html", "settings.html"))
}

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/login", getLoginHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", postLoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/settings/profile", getProfileHandler).Methods(http.MethodGet)
	r.HandleFunc("/settings/profile", postProfileHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", postLogoutHandler).Methods(http.MethodPost)
	protected := csrf.Protect([]byte("keep-it-secret-keep-it-safe----a"), csrf.Secure(false))(r)
	log.Fatal(http.ListenAndServe(":8085", protected))

}

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	d := struct {
		CSRF template.HTML
		User *user
	}{
		CSRF: csrf.TemplateField(r),
		User: nil,
	}
	if err := loginTemplate.Execute(w, d); err != nil {
		panic(err)
	}
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Println(username, password)
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	u := user{
		Name: "john",
		Age:  35,
	}
	session.Values[sessionValueKeyUser] = u
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	u := session.Values[sessionValueKeyUser].(*user)
	data := struct {
		CSRF template.HTML
		User *user
	}{
		User: u,
		CSRF: csrf.TemplateField(r),
	}
	if err := settingsTemplate.Execute(w, data); err != nil {
		panic(err)
	}
}

func postProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	u := session.Values[sessionValueKeyUser].(*user)
	name := r.FormValue("name")
	u.Name = name
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		panic(err)
	}
	u.Age = age
	session.Values[sessionValueKeyUser] = u
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
}

func postLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, sessionName)
	if err != nil {
		panic(err)
	}
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
	r.URL.Path = "/login"
	http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
}
