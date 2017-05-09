package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var (
	loginTemplate    *template.Template
	settingsTemplate *template.Template
)

func init() {
	loginTemplate = template.Must(template.ParseFiles("login.html"))
	settingsTemplate = template.Must(template.ParseFiles("settings.html"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", getLoginHandler).Methods(http.MethodGet)
	r.HandleFunc("/login", postLoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/settings", getSettingsHandler).Methods(http.MethodGet)
	r.HandleFunc("/settings", postSettingsHandler).Methods(http.MethodPost)
	r.HandleFunc("/logout", logoutHandler).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", csrf.Protect([]byte("32-byte-long-auth-key"))(r)))

}

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	loginTemplate.Execute(w, nil)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {

}

func getSettingsHandler(w http.ResponseWriter, r *http.Request) {
	settingsTemplate.Execute(w, nil)
}

func postSettingsHandler(w http.ResponseWriter, r *http.Request) {

}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

}
