package core

import (
	"bytes"
	"testing"

	"github.com/google/uuid"

	"github.com/alicebob/verssion/internal"
)

func TestPostgresDB(t *testing.T) {
	c := internal.TestDB(t)
	p := NewPGX(c)
	InterfaceTestDB(t, p)
}

func TestPostgresCurated(t *testing.T) {
	uuid.SetRand(bytes.NewBufferString("TestPostgresCuratedLongerStringfoobarbaz"))

	c := internal.TestDB(t)
	p := NewPGX(c)
	InterfaceTestCurated(t, p)
}
