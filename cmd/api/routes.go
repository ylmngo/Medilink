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
	router.HandleFunc("/scans", app.checkUser(app.scanHandler))
	// router.HandleFunc("/scans/{cat}", app.checkUser(app.ScanCatHandler))
	router.HandleFunc("/scans/upload", app.checkUser(app.scanUploadHandler))
	router.HandleFunc("/scans/upload/{cat}", app.checkUser(app.scanCatUploadHandler))
	router.HandleFunc("/reminders", app.reminderHandler)
	router.HandleFunc("/reminders/set", app.setReminderHandler)
	router.HandleFunc("/healthcheck", app.healthcheckHandler)

	router.HandleFunc("/tryupload", app.tryUploadHandler)
	router.HandleFunc("/try", app.tryHandler)

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	return router
}
