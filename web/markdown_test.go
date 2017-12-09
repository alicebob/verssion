package web

import (
	"html/template"
	"reflect"
	"testing"
)

func TestMarkdownParse(t *testing.T) {
	type cas struct {
		Src    string
		Tokens []interface{}
	}
	for i, c := range []cas{
		{
			Src:    "foo bar",
			Tokens: []interface{}{"foo bar"},
		},
		{
			Src: "[bar](link)",
			Tokens: []interface{}{
				link{href: "link", title: "bar", raw: "[bar](link)"},
			},
		},
		{
			Src: "foo [bar](link)",
			Tokens: []interface{}{
				"foo ",
				link{href: "link", title: "bar", raw: "[bar](link)"},
			},
		},
		{
			Src: "foo [bar](link) [baz title](another link)",
			Tokens: []interface{}{
				"foo ",
				link{href: "link", title: "bar", raw: "[bar](link)"},
				" ",
				link{href: "another link", title: "baz title", raw: "[baz title](another link)"},
			},
		},
		{
			Src:    "foo [bar]",
			Tokens: []interface{}{"foo [bar]"},
		},
		{
			Src:    "foo [bar](ba",
			Tokens: []interface{}{"foo [bar](ba"},
		},
		{
			Src: "foo [ [bar](link)",
			Tokens: []interface{}{
				"foo [ ",
				link{href: "link", title: "bar", raw: "[bar](link)"},
			},
		},
		{
			Src:    "foo [bar](foo",
			Tokens: []interface{}{"foo [bar](foo"},
		},
		{
			Src: "foo [foo](link)[bar](link)",
			Tokens: []interface{}{
				"foo ",
				link{href: "link", title: "foo", raw: "[foo](link)"},
				link{href: "link", title: "bar", raw: "[bar](link)"},
			},
		},
		{
			Src: "[](link)",
			Tokens: []interface{}{
				"[](link)",
			},
		},
		{
			Src: "[foo]()",
			Tokens: []interface{}{
				"[foo]()",
			},
		},
	} {
		if have, want := parseMarkdown(c.Src), c.Tokens; !reflect.DeepEqual(have, want) {
			t.Errorf("case %d: have %q, want %q", i, have, want)
		}
	}
}

func TestMarkdown(t *testing.T) {
	type cas struct {
		Src  string
		HTML string
	}
	for i, c := range []cas{
		{
			Src:  "foo bar",
			HTML: "foo bar",
		},
		{
			Src:  "foo <b>bar",
			HTML: "foo &lt;b&gt;bar",
		},
		{
			Src:  "foo [foo](http://bar)",
			HTML: `foo <a href="http://bar">foo</a>`,
		},
		{
			Src:  "foo [more words!!](http://bar)",
			HTML: `foo <a href="http://bar">more words!!</a>`,
		},
		{
			Src:  "foo [foo](http://foo)[bar](http://bar/foo/etc.html)",
			HTML: `foo <a href="http://foo">foo</a><a href="http://bar/foo/etc.html">bar</a>`,
		},
		{
			Src:  "foo [foo](mailto://bar)",
			HTML: "foo [foo](mailto://bar)",
		},
		{
			Src:  "foo [<b>foo!](http://bar)",
			HTML: `foo <a href="http://bar">&lt;b&gt;foo!</a>`,
		},
		{
			Src:  "[mariadb.org](https://mariadb.org/), [mariadb.com](https://mariadb.com/)",
			HTML: `<a href="https://mariadb.org/">mariadb.org</a>, <a href="https://mariadb.com/">mariadb.com</a>`,
		},
	} {
		if have, want := basicMarkdown(c.Src), template.HTML(c.HTML); have != want {
			t.Errorf("case %d: have %q, want %q", i, have, want)
		}
	}
}
