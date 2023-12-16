package main

import (
	"fmt"
	"lp3/internal/data"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func (app *application) scanHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) scanUploadHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/hello.html")
	if err != nil {
		app.logger.Printf("Unable to get template File: hello.html: %v\n", err)
		return
	}
	tmpl.Execute(w, nil)
}

func (app *application) scanCatUploadHandler(w http.ResponseWriter, r *http.Request) {
	cat := mux.Vars(r)["cat"]
	if r.Method == "GET" {
		tmpl, _ := template.ParseFiles("templates/target.html")
		path := fmt.Sprintf("http://localhost:8000/scans/upload/%s", cat)
		tmpl.Execute(w, path)
	} else {
		userId, _ := app.getUserFromCookie(r)

		if err := r.ParseMultipartForm(http.DefaultMaxHeaderBytes); err != nil {
			app.logger.Printf("Error while parsing multipart form: %v\n", err)
			return
		}

		files := r.MultipartForm.File["files[]"]
		for _, header := range files {
			file, err := header.Open()
			if err != nil {
				fmt.Printf("Unable to open file: %v\n", err)
				return
			}
			defer file.Close()

			f := data.NewFile(file, header, cat, userId.Id, data.SCAN)
			app.model.Fm.Insert(f)

			err = app.zip(f.Name, file)
			if err != nil {
				app.logger.Printf("Unable to zip file: %s, %v\n", err.Error(), err)
				return
			}
		}

		http.Redirect(w, r, "http://localhost:8000/", http.StatusFound)
	}
}
