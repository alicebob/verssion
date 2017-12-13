package core

import (
	"fmt"
)

type Spider interface {
	// Get the named page
	Spider(page string) (*Page, error)
}

type ErrRedirect struct {
	Page, To string
}

func (e ErrRedirect) Error() string {
	return fmt.Sprintf("%q: see page %q", e.Page, e.To)
}

type ErrNotFound struct {
	Page string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%q: no such page", e.Page)
}

type ErrNoVersion struct {
	Page string
}

func (e ErrNoVersion) Error() string {
	return fmt.Sprintf("%q: no version found", e.Page)
}
