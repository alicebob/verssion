package w

import (
	"time"
)

type Entry struct {
	Page          string
	T             time.Time
	StableVersion string
}

type DB interface {
	Recent() ([]Entry, error)
	History(...string) ([]Entry, error)
	Store(Entry) error
}
