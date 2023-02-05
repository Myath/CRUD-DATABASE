package handler

import (
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (h Handler) StudentStore(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	student := Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Fatal(err)
	}

	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		pharseCreateStudent(w, student)
		return
	}

	if err := student.storeAlreadyExistValidate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		pharseCreateStudent(w, student)
		return
	}

	grade, gpa := Grade(student.English, student.Bangla, student.Mathematics)

	student.Grade = grade
	student.GPA = gpa

	const insertStudentQuery = `
	INSERT INTO students(
		name,
		email,
		roll,
		english,
		bangla,		
		mathematics,
		grade,
		gpa 
		)  
	VALUES(
		:name,
		:email,
		:roll,
		:english,
		:bangla,		
		:mathematics,
	    :grade,
	    :gpa
		)RETURNING id;
	`
	stmt, err := h.db.PrepareNamed(insertStudentQuery)
	if err != nil {
		log.Fatal(err)
	}

	var uID int
	if err := stmt.Get(&uID, student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}
