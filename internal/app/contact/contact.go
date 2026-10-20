package contact

import (
	"net/mail"
	"strings"
)

type Contact struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
}

func (c Contact) Name() string {
	return strings.TrimSpace(c.FirstName + " " + c.LastName)
}

func (c Contact) Validate() map[string]string {
	errs := map[string]string{}
	if strings.TrimSpace(c.FirstName) == "" {
		errs["first_name"] = "First name is required."
	}
	if strings.TrimSpace(c.LastName) == "" {
		errs["last_name"] = "Last name is required."
	}
	if err := ValidateEmail(c.Email); err != "" {
		errs["email"] = err
	}
	return errs
}

func ValidateEmail(email string) string {
	email = strings.TrimSpace(email)
	if email == "" {
		return "Email is required."
	}
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "Enter a valid email address."
	}
	return ""
}

func (c Contact) matches(term string) bool {
	term = strings.ToLower(strings.TrimSpace(term))
	if term == "" {
		return true
	}
	return strings.Contains(strings.ToLower(c.Name()), term) ||
		strings.Contains(strings.ToLower(c.Email), term)
}
