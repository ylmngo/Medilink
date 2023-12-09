package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := app.getUserFromCookie(r)
	if err != nil {
		app.logger.Printf("Unable to retrieve user from cookie: %v\n", err)
		return
	}
	healthData := &map[string]string{
		"Status":      "Available",
		"Environment": app.cfg.env,
	}
	if err := app.writeJSON(w, healthData); err != nil {
		app.logger.Printf("Could not convert health data to JSON: %v\n", err)
		return
	}
}

func (app *application) indexPageHandler(w http.ResponseWriter, r *http.Request) {
	// _, err := app.getUserFromCookie(r)
	// if err == nil {
	// fmt.Println(user.Email)
	// fmt.Println(user.Name)
	// fmt.Println(user.Id)
	// }

	tmpl, err := template.ParseFiles("templates/build.html")
	if err != nil {
		app.logger.Printf("Unable to locate template file: %v\n", err)
		return
	}
	tmpl.Execute(w, nil)
}

// Login/signup handler -> rename it to something better
func (app *application) authHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		app.logger.Printf("Unable to locate template file: %v\n", err)
		return
	}
	tmpl.Execute(w, nil)
}

// Middleware function, checks if the user id is present in the cookie
func (app *application) checkUser(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := app.getUserFromCookie(r)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusNotFound)
		}
		handler.ServeHTTP(w, r)
	})
}

// Read file from disk and write to browser
func (app *application) fileViewer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fileId, _ := strconv.Atoi(id)

	file, err := app.model.Fm.GetFileById(fileId)
	if err != nil {
		app.logger.Printf("Unable to get file by id: %v\n", err)
		return
	}

	bytes, err := os.ReadFile(fmt.Sprintf("uploads/%s%s", file.Name, file.Extension))
	if err != nil {
		app.logger.Printf("Unable to get file from disk: %v\n", err)
		return
	}

	contentType := http.DetectContentType(bytes)

	w.Header().Set("Content-Type", contentType)
	w.Write(bytes)
}
