package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

var templates *template.Template

func main() {
	fmt.Println("Listening on http://localhost:8080/")
	templates = template.Must(template.ParseGlob("templates/*")) // Load all templates from "templates" folder
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/asciiart", Post) 
	//staticFiles := http.FileServer(http.Dir("./templates"))
	//http.Handle("/", staticFiles)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
