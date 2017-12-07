package core

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrCuratedNotFound = errors.New("curated ID not found")

type Page struct {
	Page          string
	T             time.Time
	StableVersion string
	Homepage      string
}

type Curated struct {
	ID          string
	CustomTitle string
	Created     time.Time
	LastUsed    time.Time
	LastUpdated time.Time
	Pages       []string
}

func (c *Curated) Title() string {
	if t := c.CustomTitle; t != "" {
		return t
	}
	return c.DefaultTitle()
}

func (c *Curated) DefaultTitle() string {
	var (
		p        = c.Pages
		ellipses = 0
		maxEls   = 4
	)
	if len(p) == 0 {
		return "[untitled feed]"
	}
	if len(p) > maxEls {
		ellipses = len(p) - maxEls
		p = p[:maxEls]
	}
	t := strings.Join(Titles(p), ", ")
	if ellipses > 0 {
		t += fmt.Sprintf(", ... (%d more)", ellipses)
	}
	return t
}

type DB interface {
	Last(string) (*Page, error) // Last spider
	Recent(int) ([]Page, error)
	CurrentAll() ([]Page, error)
	Current(...string) ([]Page, error)
	History(...string) ([]Page, error) // Newest first
	Store(Page) error
	Known() ([]string, error)

	CreateCurated() (string, error)
	LoadCurated(string) (*Curated, error) // will return (nil, nil) on not found
	CuratedSetPages(string, []string) error
	CuratedSetUsed(string) error
	CuratedSetTitle(string, string) error
}
