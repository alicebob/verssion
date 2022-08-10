package core

import (
	"context"
	"os"
	"testing"

	"github.com/alicebob/minipg"
)

func initdb(t *testing.T, addr string) DB {
	f, err := os.ReadFile("../tables.sql")
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewPostgres(addr)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := p.conn.Exec(context.Background(), string(f)); err != nil {
		t.Fatal(err)
	}

	return p
}

func TestPostgresDB(t *testing.T) {
	db := minipg.RunT(t)
	p := initdb(t, db.URL())
	t.Skip("WIP")
	InterfaceTestDB(t, p)
}

func TestPostgresCurated(t *testing.T) {
	db := minipg.RunT(t)
	p := initdb(t, db.URL())
	InterfaceTestCurated(t, p)
}
