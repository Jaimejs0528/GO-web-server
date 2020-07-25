package auth

import (
	"fmt"
	"net/http"

	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/session"
	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/template"
	"github.com/Jaimejs0528/practice/golang-udemy/cookie-exercises/my-own/controller/user"
)

// PrivateHandler route that needs a active session
type PrivateHandler struct {
	handle http.HandlerFunc
}

// NewPrivateHandler creates a new private Handler that needs a active session
func NewPrivateHandler(handle http.HandlerFunc) PrivateHandler {
	return PrivateHandler{handle}
}

func (pr PrivateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sessionCookie := getSessionCookie(r)
	if sessionCookie == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	currentSession := session.GetSession(sessionCookie.Value)
	if currentSession == nil {
		sessionCookie.MaxAge = -1
		session.RemoveSession(sessionCookie.Value)
		http.SetCookie(w, sessionCookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	sessionCookie.MaxAge = session.MaxSeconds
	session.RenovateSessionTime(sessionCookie.Value)
	http.SetCookie(w, sessionCookie)
	pr.handle.ServeHTTP(w, r)
}

// GetCurrentSession verifies if exist a current session logged
func GetCurrentSession(r *http.Request) *session.Session {
	c, err := r.Cookie("session")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	activeSession := session.GetSession(c.Value)
	return activeSession
}

func getSessionCookie(r *http.Request) *http.Cookie {
	c, err := r.Cookie("session")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return c
}

// LogOut logout the current user
func LogOut(w http.ResponseWriter, r *http.Request) {
	sessionCookie := getSessionCookie(r)
	if sessionCookie != nil {
		sessionCookie.MaxAge = -1
		http.SetCookie(w, sessionCookie)
		session.RemoveSession(sessionCookie.Value)
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

// SignIn signs in a user setting its uuid in cookies
func SignIn(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		email := r.FormValue("email")
		password := r.FormValue("password")

		u := user.GetUser(email)
		if u == nil {
			template.ServeTemplate("login.gohtml", "User doesn't exist").ServeHTTP(w, r)
			return
		}

		if !u.PasswordMatching([]byte(password)) {
			template.ServeTemplate("login.gohtml", "Incorrect email/Password").ServeHTTP(w, r)
			return
		}

		sessionUUID, err := session.CreateSession(email)
		if err != nil {
			fmt.Println(err)
			template.ServeTemplate("login.gohtml", "Ops! Something goes wrong").ServeHTTP(w, r)
		}
		c = &http.Cookie{
			Name:  "session",
			Value: sessionUUID,
		}
		c.MaxAge = session.MaxSeconds

		http.SetCookie(w, c)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session.RemoveSession(c.Value)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
