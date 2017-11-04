package w

import (
	"strings"
	"time"
)

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
	if len(c.Pages) == 0 {
		return "[untitled list]"
	}
	return strings.Join(Titles(c.Pages), ", ")
}

type DB interface {
	// Last spider
	Last(string) (*Page, error)
	Recent(int) ([]Page, error)
	// Most recent changed version
	Current(...string) ([]Page, error)
	History(...string) ([]Page, error)
	Store(Page) error
	Known() ([]string, error)

	CreateCurated() (string, error)
	LoadCurated(string) (*Curated, error)
	CuratedPages(string, []string) error
	CuratedUsed(string) error
	CuratedTitle(string, string) error
}
