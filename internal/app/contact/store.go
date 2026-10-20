package contact

import (
	"context"
	"crypto/rand"
	"errors"
	"slices"
	"strings"
	"sync"
)

var (
	ErrNotFound       = errors.New("contact not found")
	ErrDuplicateEmail = errors.New("a contact with this email already exists")
)

const PageSize = 10

type Store struct {
	mu    sync.RWMutex
	items []Contact
}

func NewStore(seed []Contact) *Store {
	s := &Store{items: slices.Clone(seed)}
	s.sort()
	return s
}

func (s *Store) sort() {
	slices.SortFunc(s.items, func(a, b Contact) int {
		if c := strings.Compare(strings.ToLower(a.LastName), strings.ToLower(b.LastName)); c != 0 {
			return c
		}
		return strings.Compare(strings.ToLower(a.FirstName), strings.ToLower(b.FirstName))
	})
}

type Page struct {
	Contacts []Contact
	Number   int  // 1-based page number
	HasMore  bool // whether a next page exists
	Total    int  // total matches across all pages
}

func (s *Store) Search(ctx context.Context, term string, page int) Page {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if page < 1 {
		page = 1
	}
	var matches []Contact
	for _, c := range s.items {
		if c.matches(term) {
			matches = append(matches, c)
		}
	}
	start := min((page-1)*PageSize, len(matches))
	end := min(start+PageSize, len(matches))
	return Page{
		Contacts: slices.Clone(matches[start:end]),
		Number:   page,
		HasMore:  end < len(matches),
		Total:    len(matches),
	}
}

func (s *Store) All(ctx context.Context) []Contact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return slices.Clone(s.items)
}

func (s *Store) Get(ctx context.Context, id string) (Contact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i := s.index(id)
	if i < 0 {
		return Contact{}, ErrNotFound
	}
	return s.items[i], nil
}

func (s *Store) Create(ctx context.Context, c Contact) (Contact, error) {
	if errs := c.Validate(); len(errs) > 0 {
		return Contact{}, &ValidationError{Fields: errs}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.emailTaken(c.Email, "") {
		return Contact{}, ErrDuplicateEmail
	}
	c.ID = rand.Text()
	c.Active = true
	s.items = append(s.items, c)
	s.sort()
	return c, nil
}

func (s *Store) Update(ctx context.Context, c Contact) (Contact, error) {
	if errs := c.Validate(); len(errs) > 0 {
		return Contact{}, &ValidationError{Fields: errs}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(c.ID)
	if i < 0 {
		return Contact{}, ErrNotFound
	}
	if s.emailTaken(c.Email, c.ID) {
		return Contact{}, ErrDuplicateEmail
	}
	c.Active = s.items[i].Active
	s.items[i] = c
	s.sort()
	return c, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(id)
	if i < 0 {
		return ErrNotFound
	}
	s.items = slices.Delete(s.items, i, i+1)
	return nil
}

func (s *Store) SetActive(ctx context.Context, ids []string, active bool) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	changed := 0
	for i := range s.items {
		if slices.Contains(ids, s.items[i].ID) && s.items[i].Active != active {
			s.items[i].Active = active
			changed++
		}
	}
	return changed
}

func (s *Store) EmailAvailable(ctx context.Context, email, exceptID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return !s.emailTaken(email, exceptID)
}

func (s *Store) index(id string) int {
	return slices.IndexFunc(s.items, func(c Contact) bool { return c.ID == id })
}

func (s *Store) emailTaken(email, exceptID string) bool {
	email = strings.ToLower(strings.TrimSpace(email))
	for _, c := range s.items {
		if c.ID != exceptID && strings.ToLower(c.Email) == email {
			return true
		}
	}
	return false
}

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	msgs := make([]string, 0, len(e.Fields))
	for _, m := range e.Fields {
		msgs = append(msgs, m)
	}
	slices.Sort(msgs)
	return strings.Join(msgs, " ")
}
