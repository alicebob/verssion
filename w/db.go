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
	Store(Entry) error
}
