package contacts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofs-cli/template/internal/app/contact"
	"github.com/gofs-cli/template/internal/ui"
)

func Index(store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.Render(w, "contacts/index", newIndexData(r, store))
	})
}

func New() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.Render(w, "contacts/new", newData{Errors: map[string]string{}})
	})
}

func Create(store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := contact.Contact{
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			Email:     r.PostFormValue("email"),
		}
		_, err := store.Create(r.Context(), c)
		if err != nil {
			d := newData{Contact: c, Errors: map[string]string{}}
			var verr *contact.ValidationError
			if errors.As(err, &verr) {
				d.Errors = verr.Fields
			} else {
				d.Errors["email"] = err.Error()
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			ui.Render(w, "contacts/new", d)
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}

// Update saves one row. On success it returns to the list it came from; on a
// validation error it renders that list again with the row open and its
// errors.
func Update(store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := contact.Contact{
			ID:        r.PathValue("id"),
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			Email:     r.PostFormValue("email"),
		}
		_, err := store.Update(r.Context(), c)
		q := r.PostFormValue("q")
		page, _ := strconv.Atoi(r.PostFormValue("page"))
		if errors.Is(err, contact.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err == nil {
			http.Redirect(w, r, listURL(q, page), http.StatusSeeOther)
			return
		}

		fieldErrors := map[string]string{}
		var verr *contact.ValidationError
		if errors.As(err, &verr) {
			fieldErrors = verr.Fields
		} else {
			fieldErrors["email"] = err.Error()
		}

		d := newIndexData(r, store)
		edited := rowData{Contact: c, Errors: fieldErrors, Editing: true, Query: q, Page: page}
		replaced := false
		for i := range d.Rows {
			if d.Rows[i].Contact.ID == c.ID {
				d.Rows[i] = edited
				replaced = true
				break
			}
		}
		if !replaced {
			d.Rows = append([]rowData{edited}, d.Rows...)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		ui.Render(w, "contacts/index", d)
	})
}

func Delete(store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := store.Delete(r.Context(), r.PathValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}
