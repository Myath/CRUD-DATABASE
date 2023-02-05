package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/justinas/nosurf"
)

type LoginAdmin struct {
	ID        int       `db:"id" form:"-"`
	Username  string    `db:"username" form:"username"`
	Password  string    `db:"password" form:"password"`
	CreatedAt time.Time `db:"created_at" form:"-"`
	CSRFToken string    `db:"-" form:"csrf_token"`
	FormError map[string]error
}

func (h Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	pareseLoginTemplate(w, LoginAdmin{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) AdminLoginProcess(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	admin := LoginAdmin{}
	if err := h.decoder.Decode(&admin, r.PostForm); err != nil {
		log.Fatal(err)
	}

	if err := admin.adminValidation(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			admin.FormError = vErr
		}
		pareseLoginTemplate(w, LoginAdmin{
			CSRFToken: nosurf.Token(r),
			FormError: admin.FormError,
		})
		return
	}

	//For isAdmin Validations

	// go h.isAdminValidate(w, r, admin.Username, admin.Password)

	var isAdmin LoginAdmin

	isAdminQuery := `SELECT * FROM admin WHERE username = $1 AND password = $2`

	if err := h.db.Get(&isAdmin, isAdminQuery, admin.Username, admin.Password); err == sql.ErrNoRows {
		pareseLoginTemplate(w, LoginAdmin{Username: admin.Username, CSRFToken: nosurf.Token(r), FormError: map[string]error{
			"Username": fmt.Errorf("The username/password doesn't match."),
		}})
		return
	}

	h.sessionManager.Put(r.Context(), "username", admin.Username)

	http.Redirect(w, r, "/student/list", http.StatusSeeOther)

}

// func (h Handler) isAdminValidate(w http.ResponseWriter, r *http.Request, username string, password string){

// 	var isAdmin LoginAdmin

// 	isAdminQuery := `SELECT * FROM admin WHERE username = $1 AND password = $2`

// 	if err := h.db.Get(&isAdmin, isAdminQuery, username, password); err != nil {
// 		pareseLoginTemplate(w, LoginAdmin{Username: username,  CSRFToken: nosurf.Token(r), FormError: map[string]error{
// 			"Username": fmt.Errorf("The username/password doesn't match."),
// 		}})
// 		return
// 	}
// 	return
// }

func pareseLoginTemplate(w http.ResponseWriter, data any) {
	t, err := template.ParseFiles("templates/header.html", "templates/footer.html", "templates/adminlogin.html")
	if err != nil {
		log.Fatalf("%v", err)
	}
	t.ExecuteTemplate(w, "adminlogin.html", data)
}
