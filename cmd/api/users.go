package main

import (
	"database/sql"
	"errors"
	"fmt"
	"lp3/internal/data"
	"net/http"
)

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	user := data.NewUser(username, password, email)
	if err := app.model.Um.Insert(user); err != nil {
		app.logger.Printf("Unable to insert user to DB: %v\n", err)
		return
	}

	http.Redirect(w, r, "/auth", http.StatusFound)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := app.model.Um.GetUserByEmail(email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			app.logger.Printf("Invalid Email Address: %v\n", err)
			http.Redirect(w, r, "../healthcheck", http.StatusFound)
		default:
			app.logger.Printf("Unable to retrieve user: %v\n", err)
			http.Redirect(w, r, "../healthcheck", http.StatusFound)
		}
	}
	if !user.CompareHashAndPassword(password) {
		app.logger.Println("Incorrect Password, Please try again")
		http.Redirect(w, r, "../healthcheck", http.StatusFound)
	}
	cookie := &http.Cookie{
		Name:  "user-cookie",
		Value: fmt.Sprintf("%d", user.Id),
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "../", http.StatusFound)
}

// TODO: send email or sms to user's email or phone number, reminding them
func (app *application) reminderHandler(w http.ResponseWriter, r *http.Request) {
	user, err := app.getUserFromCookie(r)
	if err != nil {
		app.logger.Printf("Unable to retrieve user from cookie: %v\n", err)
		return
	}
	err = app.mailer.Send(user.Email, user.Name)
	if err != nil {
		app.logger.Printf("Unable to send mail to user: %v\n", err)
		return
	}
}
