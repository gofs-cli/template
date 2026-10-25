package contacts

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gofs-cli/template/internal/app/contact"
)

type newData struct {
	Contact contact.Contact
	Errors  map[string]string
}

// rowData is one contact and its edit form. Editing is true when the server
// re-renders a submitted row with validation errors.
type rowData struct {
	Contact contact.Contact
	Errors  map[string]string
	Editing bool
	Query   string
	Page    int
}

type indexData struct {
	Query   string
	Page    contact.Page
	Rows    []rowData
	MoreURL string
}

func newIndexData(r *http.Request, store *contact.Store) indexData {
	q := r.FormValue("q")
	page, _ := strconv.Atoi(r.FormValue("page"))
	result := store.Search(r.Context(), q, page)
	d := indexData{Query: q, Page: result}
	if result.HasMore {
		d.MoreURL = listURL(q, result.Number+1)
	}
	for _, c := range result.Contacts {
		d.Rows = append(d.Rows, rowData{Contact: c, Query: q, Page: result.Number})
	}
	return d
}

// listURL builds the contacts URL for a search and page, so a form post can
// return to the list it came from.
func listURL(q string, page int) string {
	values := url.Values{}
	if q != "" {
		values.Set("q", q)
	}
	if page > 1 {
		values.Set("page", strconv.Itoa(page))
	}
	if len(values) == 0 {
		return "/contacts"
	}
	return "/contacts?" + values.Encode()
}
