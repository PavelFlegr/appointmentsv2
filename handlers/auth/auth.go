package auth

import (
	"appointmentsv2/global"
	"appointmentsv2/services/auth"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	user, err := auth.GetUserByEmail(email)
	if err != nil {
		_, err := fmt.Fprintln(w, "Account with this email does not exist")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		_, err := fmt.Fprintln(w, "Password is incorrect")
		if err != nil {
			log.Println(err)
		}
		return
	}

	cookie, _ := global.Sc.Encode("userId", user.Id)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Value:    cookie,
		Name:     "userId",
		SameSite: http.SameSiteStrictMode},
	)

	w.Header().Set("HX-Redirect", "/manage/appointment")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "userId",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {
	email, err := mail.ParseAddress(r.PostFormValue("email"))
	if err != nil {
		_, err := fmt.Fprintln(w, "Invalid email address")
		if err != nil {
			log.Println(err)
		}
		return
	}

	password := r.PostFormValue("password")
	password2 := r.PostFormValue("password2")

	if len(password) < 5 {
		_, err := fmt.Fprintln(w, "Password must be at least 5 characters long")
		if err != nil {
			log.Println(err)
		}
		return
	}

	if password != password2 {
		_, err := fmt.Fprintln(w, "Passwords don't match")
		if err != nil {
			log.Println(err)
		}
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user := auth.User{
		Email:    email.Address,
		Password: string(hash),
	}
	err = auth.CreateUser(&user)
	if err != nil {
		_, err = fmt.Fprintln(w, "An account with this email already exists")
		if err != nil {
			log.Println(err)
		}
	}

	cookie, _ := global.Sc.Encode("userId", user.Id)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Value:    cookie,
		Name:     "userId",
		SameSite: http.SameSiteStrictMode},
	)

	w.Header().Set("HX-Redirect", "/manage/appointment")
}
