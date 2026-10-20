package contacts

import (
	"errors"
	"net/http"

	"github.com/gofs-cli/template/internal/app/contact"
	"github.com/gofs-cli/template/internal/ui"
)

func Index(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		d := newIndexData(r, store)
		switch {
		case ui.WantsFragment(r) && d.Page.Number > 1:
			ui.RenderFragment(w, "contacts/page", d)
		case ui.WantsFragment(r):
			ui.RenderFragment(w, "contacts/results", d)
		default:
			ui.RenderLayoutContent(w, "layouts/main", "contacts/index", d)
		}
	})
}

func Row(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := store.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if ui.WantsFragment(r) {
			ui.RenderFragment(w, "contacts/row", rowData{Contact: c})
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}

func EditRow(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := store.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		d := editData{Contact: c}
		if ui.WantsFragment(r) {
			ui.RenderFragment(w, "contacts/edit-row", rowData{Contact: c})
			return
		}
		ui.RenderLayoutContent(w, "layouts/main", "contacts/edit", d)
	})
}

func Update(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := contact.Contact{
			ID:        r.PathValue("id"),
			FirstName: r.PostFormValue("first_name"),
			LastName:  r.PostFormValue("last_name"),
			Email:     r.PostFormValue("email"),
		}
		saved, err := store.Update(r.Context(), c)
		if errors.Is(err, contact.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			if ui.WantsFragment(r) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				ui.RenderFragment(w, "contacts/edit-row", rowData{Contact: c})
				return
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			ui.RenderLayoutContent(w, "layouts/main", "contacts/edit", editData{Contact: c, Error: err.Error()})
			return
		}
		if ui.WantsFragment(r) {
			ui.RenderFragment(w, "contacts/row", rowData{Contact: saved})
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}

func Delete(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := store.Delete(r.Context(), r.PathValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if ui.WantsFragment(r) {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}

func New(ui *ui.UI) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ui.RenderLayoutContent(w, "layouts/main", "contacts/new", newData{Errors: map[string]string{}})
	})
}

func Create(ui *ui.UI, store *contact.Store) http.Handler {
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
			ui.RenderLayoutContent(w, "layouts/main", "contacts/new", d)
			return
		}
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	})
}

func Validate(ui *ui.UI, store *contact.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		msg := ""
		switch {
		case r.PostForm.Has("email"):
			email := r.PostFormValue("email")
			msg = contact.ValidateEmail(email)
			if msg == "" && !store.EmailAvailable(r.Context(), email, r.PostFormValue("id")) {
				msg = contact.ErrDuplicateEmail.Error()
			}
		default:
			http.Error(w, "unknown field", http.StatusBadRequest)
			return
		}
		ui.RenderFragment(w, "contacts/field-error", struct{ Message string }{msg})
	})
}
