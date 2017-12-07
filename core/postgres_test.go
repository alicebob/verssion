// +build integration

package core

import (
	"testing"
)

func TestPostgresDB(t *testing.T) {
	p, err := NewPostgres("postgresql:///verssion")
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"page"} {
		if _, err := p.conn.Exec("DELETE FROM " + table); err != nil {
			t.Fatal(err)
		}
	}

	InterfaceTestDB(t, p)
}
