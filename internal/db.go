package internal

import (
	"context"
	"testing"

	"github.com/alicebob/pgsnap"
	"github.com/jackc/pgx/v4/pgxpool"
)

var tables = []string{"page", "curated", "curated_pages"}

// return the connection to a cleaned, setup postgres db. See also pgsnap docs (and PGPROXY).
func TestDB(t *testing.T) *pgxpool.Pool {
	ctx := context.Background()
	c := pgsnap.RunEnvPGXPool(t)

	for _, table := range tables {
		if _, err := c.Exec(ctx, "DELETE FROM "+table); err != nil {
			t.Fatal(err)
		}
	}

	return c
}
