package core

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

func TestCleanSpace(t *testing.T) {
	for orig, want := range map[string]string{
		" foo   ":          "foo",
		" f   oo   ":       "f oo",
		"new\nline":        "new\nline",
		"new  \n line":     "new\nline",
		"  new  \n\n line": "new\nline",
		"\n":               "",
		"\n\n":             "",
	} {
		if have := cleanSpace(orig); have != want {
			t.Errorf("have %q, want %q", have, want)
		}
	}
}

func TestFindTables(t *testing.T) {
	type cas struct {
		Html string
		Want []Table
	}
	cases := []cas{
		{
			Html: `<table>string</table>`,
			Want: []Table{
				Table{
					Rows: [][]string(nil),
				},
			},
		},
		{
			Html: `<html><body><table><tr><td>foo</td><td>bar</td></table>`,
			Want: []Table{
				{
					Rows: [][]string{{"foo", "bar"}},
				},
			},
		},
		{
			Html: `<html><body><table><tr><td><a href="https://foo">bar</a></td></table>`,
			Want: []Table{
				{
					Rows: [][]string{{"[bar](https://foo)"}},
				},
			},
		},
		{
			Html: `<html><body><table><tr><td><a href="//foo">bar</a></td></table>`,
			Want: []Table{
				{
					Rows: [][]string{{"[bar](https://foo)"}},
				},
			},
		},
		{
			Html: `<html><body><table><tr><td><a href="https://foo">[bar]</a></td></table>`,
			Want: []Table{
				{
					Rows: [][]string{{`[\[bar\]](https://foo)`}},
				},
			},
		},
		{
			Html: `<html><body><table><tr><td>no<a href="https://foo">[±]</a>link</td></table>`,
			Want: []Table{
				{
					Rows: [][]string{{"nolink"}},
				},
			},
		},
	}
	for i, c := range cases {
		doc, err := html.Parse(bytes.NewBufferString(c.Html))
		if err != nil {
			t.Fatal(err)
		}

		d, err := FindTables(doc)
		if err != nil {
			t.Fatal(err)
		}
		if have, want := d, c.Want; !reflect.DeepEqual(have, want) {
			t.Errorf("case %d: have %#v, want %#v", i, have, want)
		}
	}
}

func TestFindTablesReal(t *testing.T) {
	r, err := os.Open("./data/git.html")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	doc, err := html.Parse(r)
	if err != nil {
		t.Fatal(err)
	}

	ts, err := FindTables(doc)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(ts), 6; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
	t1 := Table{Rows: [][]string{
		[]string{""},
		[]string{"A command-line session showing repository creation, addition of a file, and remote synchronization"},
		[]string{"Original author(s)", "Linus Torvalds"},
		[]string{"Developer(s)", "Junio Hamano and others"},
		[]string{"Initial release", "7 April 2005"},
		[]string{""},
		[]string{"Stable release", "2.14.2 / 22 September 2017"},
		[]string{""},
		[]string{"Repository", "[git-scm.com/downloads](https://git-scm.com/downloads)"},
		[]string{"Development status", "Active"},
		[]string{"Written in", "C, Shell, Perl, Tcl, Python"},
		[]string{"Operating system", "POSIX: Linux, Windows, macOS"},
		[]string{"Platform", "IA-32, x86-64"},
		[]string{"Available in", "English"},
		[]string{"Type", "Version control"},
		[]string{"License", "GNU GPL v2 and GNU LGPL v2.1"},
		[]string{"Website", "[git-scm.com](https://git-scm.com)"},
	}}
	if have, want := ts[0], t1; !reflect.DeepEqual(have, want) {
		t.Errorf("have %#v\nwant %#v", have, want)
	}
}
