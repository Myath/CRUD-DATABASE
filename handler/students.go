package handler

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/justinas/nosurf"
)

type Student struct {
	ID          int              `db:"id" form:"-"`
	Name        string           `db:"name" form:"name"`
	Email       string           `db:"email" form:"email"`
	Roll        int              `db:"roll" form:"roll"`
	English     int              `db:"english" form:"eng"`
	Bangla      int              `db:"bangla" form:"ban"`
	Mathematics int              `db:"mathematics" form:"math"`
	Grade       string           `db:"grade" form:"-"`
	GPA         float64          `db:"gpa" form:"-"`
	Status      bool             `db:"status" form:"status"`
	CreatedAt   time.Time        `db:"created_at" form:"-"`
	UpdatedAt   time.Time        `db:"updated_at" form:"-"`
	DeletedAt   sql.NullTime     `db:"deleted_at" form:"-"`
	FormError   map[string]error `db:"-"`
	CSRFToken   string           `db:"-" form:"csrf_token"`
}

type StudentList struct {
	Students []Student `json:"students"`
}

func (h Handler) StudentsList(w http.ResponseWriter, r *http.Request) {
	const listQuery = `SELECT * FROM students WHERE deleted_at IS NULL`
	var student []Student

	if err := h.db.Select(&student, listQuery); err != nil {
		log.Fatal(err)
	}

	t, err := template.ParseFiles("templates/admin/student/students-list.html")
	if err != nil {
		log.Fatal(err)
	}

	sort.SliceStable(student, func(i, j int) bool {
		return student[i].ID < student[j].ID
	})

	t.Execute(w, student)
}

func (h Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	pharseCreateStudent(w, Student{
		CSRFToken: nosurf.Token(r),
	})
}
