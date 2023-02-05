package handler

import (
	"html/template"
	"log"
	"net/http"
)

func pharseCreateStudent(w http.ResponseWriter, data any) {
	t := template.New("create student")
	t = template.Must(t.ParseFiles("templates/admin/student/create-students.html", "templates/admin/student/_form.html"))

	if err := t.ExecuteTemplate(w,"create-students.html", data); err != nil {
		log.Fatal(err)
	}
}

func pharseEditStudent(w http.ResponseWriter, data any) {
	t := template.New("edit student")
	t = template.Must(t.ParseFiles("templates/admin/student/edit-student.html", "templates/admin/student/_form.html"))

	if err := t.ExecuteTemplate(w, "edit-student.html", data); err != nil{
		log.Fatal(err)
	}
}
