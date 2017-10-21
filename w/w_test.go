package w

import (
	"os"
	"testing"
)

func TestStableVersion(t *testing.T) {
	type cas struct {
		Filename string
		Version  string
	}
	cases := []cas{
		{
			Filename: "PostgreSQL.xml",
			Version:  "10.0",
		},
		{
			Filename: "MySQL.xml",
			Version:  "5.7.20",
		},
		{
			Filename: "debian.xml",
			Version:  "9.2 (Stretch)",
		},
	}

	for _, c := range cases {
		r, err := os.Open("./data/" + c.Filename)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		pt, err := Parse(r)
		if err != nil {
			t.Fatal(err)
		}
		version := StableVersion(pt)
		if have, want := version, c.Version; have != want {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
