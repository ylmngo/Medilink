package main

import (
	"fmt"
	"io"
	"lp3/internal/data"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
)

// Gets distinct prescription categories and displays them as folders
func (app *application) prescriptionHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := app.getUserFromCookie(r)
	categories, err := app.model.Fm.GetCategoriesByUser(user.Id, data.PRESCRIPTION)
	if err != nil {
		app.logger.Printf("Unable to retrieve cateogires from User: %v\n", err)
		return
	}
	tmpl, err := template.ParseFiles("templates/categories.html")
	if err != nil {
		app.logger.Printf("Unable to locate cateogires.html: %v\n", err)
		return
	}
	tmpl.Execute(w, categories)
}

// Gets files within the specific category and lists them
func (app *application) categoryHandler(w http.ResponseWriter, r *http.Request) {
	cat := mux.Vars(r)["cat"]
	files, err := app.model.Fm.GetFilesByCategory(cat, data.PRESCRIPTION)
	if err != nil {
		app.logger.Printf("Unable to get Files from category: %s, %v\n", cat, err)
		return
	}
	tmpl, err := template.ParseFiles("templates/files.html")
	if err != nil {
		app.logger.Printf("Unable to locate files.html: %v\n", err)
		return
	}

	tmpl.Execute(w, files)
}

// Handles file uploads and saves them to disk
// TODO: make it more generic for handling scans/report uploads
// TODO: File Compression
func (app *application) prescriptionUploadHandler(w http.ResponseWriter, r *http.Request) {
	user, _ := app.getUserFromCookie(r)
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles("templates/form.html")
		if err != nil {
			app.logger.Printf("Unable to locate form template: %v\n", err)
			return
		}
		tmpl.Execute(w, nil)
	} else {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			app.logger.Printf("Unable to parse multipart form: %v\n", err)
			return
		}

		formFile, fileHeader, err := r.FormFile("prescription")
		if err != nil {
			app.logger.Printf("Unable to get form file from multipart form: %v\n", err)
			return
		}
		defer formFile.Close()

		// Create a new instance of File struct and Insert new row in files table
		f := data.NewFile(formFile, fileHeader, r.FormValue("category"), user.Id, data.PRESCRIPTION)
		app.model.Fm.Insert(f)

		// Create destination File as UUID.Extension in the uploads folder
		dst, err := os.Create(fmt.Sprintf("./uploads/%s%s", f.Name, f.Extension))
		if err != nil {
			app.logger.Printf("Unable to create destination file: %v\n", err)
			return
		}
		defer dst.Close()

		if _, err = io.Copy(dst, formFile); err != nil {
			app.logger.Printf("Unable to copy data from form file to destination: %v\n", err)
			return
		}

		redirectPath := fmt.Sprintf("%s/prescriptions", r.Host)
		http.Redirect(w, r, redirectPath, http.StatusFound)
	}
}
