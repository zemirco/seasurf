package main

import (
	"html/template"
	"log"
	"net/http"

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
	loginTemplate = template.Must(template.ParseFiles("login.html"))
	settingsTemplate = template.Must(template.ParseFiles("settings.html"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", getLoginHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", postLoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/settings/profile", getProfileHandler).Methods(http.MethodGet)
	r.HandleFunc("/settings/profile", postProfileHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler).Methods(http.MethodPost)
	protected := csrf.Protect([]byte("keep-it-secret-keep-it-safe----a"), csrf.Secure(false))(r)
	log.Fatal(http.ListenAndServe(":8085", protected))

}

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	d := struct {
		CSRF template.HTML
	}{
		CSRF: csrf.TemplateField(r),
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
	u := &user{
		Name: "john",
		Age:  35,
	}
	session.Values[sessionValueKeyUser] = u
	if err := session.Save(r, w); err != nil {
		panic(err)
	}
}

func getProfileHandler(w http.ResponseWriter, r *http.Request) {
	settingsTemplate.Execute(w, nil)
}

func postProfileHandler(w http.ResponseWriter, r *http.Request) {

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

}
