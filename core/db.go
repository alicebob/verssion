package core

import (
	"errors"
	"time"
)

var ErrCuratedNotFound = errors.New("curated ID not found")

type Page struct {
	Page          string
	T             time.Time
	StableVersion string
	Homepage      string
}

type SortBy int

const (
	SpiderT SortBy = iota
	Alphabet
)

func (s SortBy) String() string {
	switch s {
	case SpiderT:
		return "timestamp DESC"
	case Alphabet:
		return "page ASC"
	default:
		panic("...")
	}
}

type DB interface {
	Last(string) (*Page, error) // Last spider
	Current(limit int, order SortBy) ([]Page, error)
	CurrentIn(...string) ([]Page, error)
	History(...string) ([]Page, error) // Newest first
	Store(Page) error
	Known() ([]string, error)

	CreateCurated() (string, error)
	LoadCurated(string) (*Curated, error) // will return (nil, nil) on not found
	CuratedSetPages(string, []string) error
	CuratedSetUsed(string) error
	CuratedSetTitle(string, string) error
}
