package core

import (
	"context"
	"testing"

	"github.com/alicebob/pgsnap"
)

var tables = []string{"page", "curated", "curated_pages"}

func initdb(t *testing.T, addr string) DB {
	p, err := NewPostgres(addr)
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range tables {
		if _, err := p.conn.Exec(context.Background(), "DELETE FROM "+table); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func TestPostgresDB(t *testing.T) {
	addr := pgsnap.RunEnv(t, "postgresql:///verssion")

	p := initdb(t, addr)
	InterfaceTestDB(t, p)
}

func TestPostgresCurated(t *testing.T) {
	addr := pgsnap.RunEnv(t, "postgresql:///verssion")

	p := initdb(t, addr)
	InterfaceTestCurated(t, p)
}
