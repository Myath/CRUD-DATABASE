package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) StudentEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	var student Student
	student.ID = uID

	const editQuery = `SELECT * FROM students WHERE id = $1 AND deleted_at IS NULL`

	if err := h.db.Get(&student, editQuery, uID); err != nil {
		log.Fatal(err)
	}

	student.CSRFToken = nosurf.Token(r)

	pharseEditStudent(w, student)
}

func (h Handler) StudentUpdate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	student := Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Fatal(err)
	}

	student.ID = uID
	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			student.FormError = vErr
		}
		pharseEditStudent(w, student)
		return
	}

	grade, gpa := Grade(student.English, student.Bangla, student.Mathematics)

	student.Grade = grade
	student.GPA = gpa

	const updateQuery = `UPDATE students
		SET name = :name, 
			email = :email,
			roll = :roll,
			english = :english,
			bangla = :bangla,
			mathematics = :mathematics,
			grade = :grade,
			gpa = :gpa,
			status = :status
		WHERE id = :id
		RETURNING id;
	`

	stmt, err := h.db.PrepareNamed(updateQuery)
	if err != nil {
		log.Fatal(err)
	}

	if err := stmt.Get(&uID, student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)
}
