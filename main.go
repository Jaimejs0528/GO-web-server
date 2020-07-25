package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/auth"
	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/session"
	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/template"
	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/user"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth.SignIn(w, r)
		return
	}
	if auth.GetCurrentSession(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	template.ServeTemplate("login.gohtml", nil).ServeHTTP(w, r)
}

func signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth.SignUp(w, r)
		return
	}
	if auth.GetCurrentSession(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	template.ServeTemplate("signup.gohtml", nil).ServeHTTP(w, r)
}

func index(w http.ResponseWriter, r *http.Request) {
	user := user.GetUser(auth.GetCurrentSession(r).UserEmail)
	template.ServeTemplate("index.gohtml", user).ServeHTTP(w, r)
}

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public", fs))
	http.Handle("/", auth.NewPrivateHandler(index))
	http.HandleFunc("/login", login)
	http.HandleFunc("/sign-up", signUp)
	http.HandleFunc("/logout", auth.LogOut)
	fmt.Println("serving on port 3000")
	go session.GoCleanExpireSessions()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
