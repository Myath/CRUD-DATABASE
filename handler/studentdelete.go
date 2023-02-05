package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (h Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	const deleteQuery = `
	DELETE FROM students where id = $1`

	res := h.db.MustExec(deleteQuery, uID)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/student/list", http.StatusSeeOther)

}
