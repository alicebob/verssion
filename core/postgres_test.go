// +build integration

package core

import (
	"testing"
)

var tables = []string{"page", "curated", "curated_pages"}

func initdb(t *testing.T) DB {
	p, err := NewPostgres("postgresql:///verssion")
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range tables {
		if _, err := p.conn.Exec("DELETE FROM " + table); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func TestPostgresDB(t *testing.T) {
	p := initdb(t)
	InterfaceTestDB(t, p)
}

func TestPostgresCurated(t *testing.T) {
	p := initdb(t)

	InterfaceTestCurated(t, p)
}
