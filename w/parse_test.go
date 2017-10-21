package w

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type cas struct {
		XML  string
		Want []Template
	}
	cases := []cas{
		{
			XML: `<root>string</root>`,
		},
		{
			XML: `<root><template><title>Use mdy dates</title></template></root>`,
			Want: []Template{
				{
					Title: "Use mdy dates",
				},
			},
		},
		{
			XML: `<root><template><title>Use mdy dates</title><part><name>date</name><equals>=</equals><value>March 2014</value></part></template></root>`,
			Want: []Template{
				{
					Title: "Use mdy dates",
					NamedParts: map[string]string{
						"date": "March 2014",
					},
				},
			},
		},
		{
			XML: `<root>
<template lineStart="1"><title>Infobox software
</title><part><name> logo                   </name><equals>=</equals><value> [[File:Postgresql elephant.svg|120px]]
</value></part><part><name> developer              </name><equals>=</equals><value> PostgreSQL Global Development Group
</value></part><part><name> released               </name><equals>=</equals><value> <template><title>Start date and age</title><part><name index="1"/><value>1996</value></part><part><name index="2"/><value>7</value></part><part><name index="3"/><value>8</value></part><part><name>br</name><equals>=</equals><value>yes</value></part><part><name>df</name><equals>=</equals><value>yes</value></part></template><ext><name>ref</name><attr> name=&quot;birthday&quot; </attr></ext>
</value></part><part><name> latest release version </name><equals>=</equals><value> <comment>&lt;!-- If you update this, remember to also update [[Comparison of relational database management systems]]--&gt;</comment>10.0
</value></part><part><name> latest release date    </name><equals>=</equals><value> <template><title>Start date and age</title><part><name index="1"/><value>2017</value></part><part><name index="2"/><value>10</value></part><part><name index="3"/><value>05</value></part><part><name>br</name><equals>=</equals><value>yes</value></part><part><name>df</name><equals>=</equals><value>yes</value></part></template><ext><name>ref</name><attr/><inner>{{cite web
 | url       = https://www.postgresql.org/about/news/1786/
 | title     = PostgreSQL 10 Released
 | publisher = The PostgreSQL Global Development Group
 | date      = 2017-10-05
 | website   = PostgreSQL
 |accessdate = 2017-10-06
}}</inner><close>&lt;/ref&gt;</close></ext>
</value></part>
</template>
</root>`,
			Want: []Template{
				{
					Title: "Infobox software",
					NamedParts: map[string]string{
						"logo":                   "[[File:Postgresql elephant.svg|120px]]",
						"developer":              "PostgreSQL Global Development Group",
						"latest release version": "10.0",
						"latest release date":    "", // nested template
						"released":               "", // nested template
					},
				},
			},
		},
	}
	for i, c := range cases {
		d, err := Parse(bytes.NewBufferString(c.XML))
		if err != nil {
			t.Fatal(err)
		}
		if have, want := d, c.Want; !reflect.DeepEqual(have, want) {
			t.Errorf("case %d: have %#v, want %#v", i, have, want)
		}
	}
}

func TestParseReal(t *testing.T) {
	r, err := os.Open("./data/PostgreSQL.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	ts, err := Parse(r)
	if err != nil {
		t.Fatal(err)
	}
	if have, want := len(ts), 82; have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}

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
