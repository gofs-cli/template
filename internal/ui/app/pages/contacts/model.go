package contacts

import (
	"net/http"
	"strconv"

	"github.com/gofs-cli/template/internal/app/contact"
)

type editData struct {
	Contact contact.Contact
	Error   string
}

type newData struct {
	Contact contact.Contact
	Errors  map[string]string
}

type rowData struct {
	Contact contact.Contact
	Error   string
}

type indexData struct {
	Query    string
	Page     contact.Page
	Rows     []rowData
	NextPage int
}

func newIndexData(r *http.Request, store *contact.Store) indexData {
	q := r.FormValue("q")
	page, _ := strconv.Atoi(r.FormValue("page"))
	result := store.Search(r.Context(), q, page)
	d := indexData{
		Query:    q,
		Page:     result,
		NextPage: result.Number + 1,
	}
	for _, c := range result.Contacts {
		d.Rows = append(d.Rows, rowData{Contact: c})
	}
	return d
}
