package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"

	logger "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

const (
	HTTP_PORT = "88"
	HOME      = "/opt/openbook/landing"
	// HOME      = "./build"
	RESOURCES = HOME + "/resources"
	DATA_FILE = HOME + "/users.csv"
)

type EmailForm struct {
	Name  string
	Email string
}

var dataFile *os.File

func index(rw http.ResponseWriter, req *http.Request) {
	file := path.Join(RESOURCES, "index.html")
	t, err := template.ParseFiles(file)
	if err != nil {
		logger.Errorf("Can't parse template file %s: %s", file, err)
	}
	t.Execute(rw, map[string]interface{}{
		"r": req.FormValue("r"),
	})
}

func save(rw http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		logger.Error(err)
	}

	form := new(EmailForm)
	decoder := schema.NewDecoder()
	if err := decoder.Decode(form, req.PostForm); err != nil {
		logger.Error(err)
	}

	if form.Email != "" {
		f, _ := os.OpenFile(DATA_FILE, os.O_APPEND|os.O_WRONLY, 0666)
		if _, err := f.WriteString(fmt.Sprintf("%s;%s\n", form.Email, form.Name)); err != nil {
			logger.Errorf("Error while writing to file: %s", err)
		}
		f.Close()
		http.Redirect(rw, req, "/?r=success", http.StatusFound)
		return
	}

	http.Redirect(rw, req, "/?r=fail", http.StatusFound)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/", save).Methods("POST")

	p := path.Join(RESOURCES, "static")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(p))))

	logger.Infof("Listening webserver on port %s", HTTP_PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%s", HTTP_PORT), router)
	if err != nil {
		panic(err)
	}
}
