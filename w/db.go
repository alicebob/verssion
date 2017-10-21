package w

import (
	"time"
)

type Entry struct {
	Page          string
	Revision      int
	T             time.Time
	StableVersion string
}

type DB interface {
	Load(string) (*Entry, error)
	Store(Entry) error
}
