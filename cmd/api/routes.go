package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) route() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", app.indexPageHandler)
	router.HandleFunc("/auth", app.authHandler)
	router.HandleFunc("/login", app.loginHandler).Methods("POST")
	router.HandleFunc("/register", app.registerHandler).Methods("POST")
	router.HandleFunc("/prescriptions/upload", app.checkUser(app.prescriptionUploadHandler))
	router.HandleFunc("/prescriptions", app.checkUser(app.prescriptionHandler))
	router.HandleFunc("/prescriptions/category/{cat}", app.checkUser(app.categoryHandler))
	router.HandleFunc("/prescriptions/category/file/{id}", app.fileViewer)
	router.HandleFunc("/reminders", app.reminderHandler)
	router.HandleFunc("/healthcheck", app.healthcheckHandler)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return router
}
