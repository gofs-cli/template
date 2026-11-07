package contact

import "fmt"

func SampleContacts() []Contact {
	names := [][2]string{
		{"Ada", "Lovelace"},
		{"Alan", "Turing"},
		{"Grace", "Hopper"},
		{"Edsger", "Dijkstra"},
		{"Barbara", "Liskov"},
		{"Donald", "Knuth"},
		{"Margaret", "Hamilton"},
		{"Dennis", "Ritchie"},
		{"Ken", "Thompson"},
		{"Rob", "Pike"},
		{"Frances", "Allen"},
		{"Tony", "Hoare"},
		{"Niklaus", "Wirth"},
		{"Adele", "Goldberg"},
		{"John", "Backus"},
		{"Radia", "Perlman"},
		{"Leslie", "Lamport"},
		{"Jean", "Bartik"},
		{"Brian", "Kernighan"},
		{"Katherine", "Johnson"},
		{"Bjarne", "Stroustrup"},
		{"Anita", "Borg"},
		{"Guido", "van Rossum"},
		{"Mary", "Kenneth Keller"},
		{"Linus", "Torvalds"},
		{"Shafi", "Goldwasser"},
	}
	contacts := make([]Contact, 0, len(names))
	for i, n := range names {
		contacts = append(contacts, Contact{
			ID:        fmt.Sprintf("c%03d", i+1),
			FirstName: n[0],
			LastName:  n[1],
			Email:     fmt.Sprintf("%s.%s@example.com", slug(n[0]), slug(n[1])),
			Active:    i%7 != 3,
		})
	}
	return contacts
}

func slug(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'A' && c <= 'Z':
			out = append(out, c+'a'-'A')
		case c >= 'a' && c <= 'z':
			out = append(out, c)
		case c == ' ':
			out = append(out, '-')
		}
	}
	return string(out)
}
