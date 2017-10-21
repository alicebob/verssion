package w

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Template struct {
	Title      string
	NamedParts map[string]string
}

// Parse returns all top-level templates
func Parse(r io.Reader) ([]Template, error) {
	d := xml.NewDecoder(r)
	root, err := d.Token()
	if err != nil {
		return nil, err
	}
	if _, ok := root.(xml.StartElement); !ok {
		return nil, fmt.Errorf("expected to start with '<root>'")
	}

	var res []Template
	for {
		t, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch tok := t.(type) {
		case xml.StartElement:
			switch tok.Name {
			case tag("template"):
				s, err := tTemplate(d)
				if err != nil {
					return nil, err
				}
				res = append(res, *s)
			default:
				// fmt.Printf("ignoring %s\n", tok.Name)
				if err := d.Skip(); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			// That should be the </root>. No need to check for EOF
			return res, nil
		}
	}
}

func tTemplate(d *xml.Decoder) (*Template, error) {
	res := &Template{}
	for {
		t, err := d.Token()
		if err != nil {
			return res, err
		}
		switch tok := t.(type) {
		case xml.StartElement:
			switch tok.Name {
			case tag("title"):
				t, err := tGen(d)
				if err != nil {
					return nil, err
				}
				res.Title = trim(t)
			case tag("part"):
				name, value, err := tPart(d)
				if err != nil {
					return nil, err
				}
				if name != "" {
					if res.NamedParts == nil {
						res.NamedParts = map[string]string{}
					}
					res.NamedParts[name] = value
				}
			default:
				// fmt.Sprintf("template ignoring: %v\n", tok.Name)
				if err := d.Skip(); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			return res, nil
		}
	}
}

func tGen(d *xml.Decoder) (string, error) {
	res := ""
	for {
		t, err := d.Token()
		if err != nil {
			return "", err
		}
		switch tok := t.(type) {
		case xml.StartElement:
			// fmt.Printf("gen ignoring open: %q\n", tok.Name)
			if err := d.Skip(); err != nil {
				return "", err
			}
		case xml.EndElement:
			return res, nil
		case xml.CharData:
			res += string(tok)
		}
	}
}

// name and value from
// <part><name>date</name><equals>=</equals><value>March 2014</value></part>
func tPart(d *xml.Decoder) (string, string, error) {
	name, value := "", ""
	seenAssignment := false
	for {
		t, err := d.Token()
		if err != nil {
			return "", "", err
		}
		switch tok := t.(type) {
		case xml.StartElement:
			switch tok.Name {
			case tag("name"):
				name, err = tGen(d)
				if err != nil {
					return "", "", err
				}
			case tag("value"):
				value, err = tGen(d)
				if err != nil {
					return "", "", err
				}
			case tag("equals"):
				seenAssignment = true
				if err := d.Skip(); err != nil {
					return "", "", err
				}
			default:
				if err := d.Skip(); err != nil {
					return "", "", err
				}
			}
		case xml.EndElement:
			if seenAssignment {
				return trim(name), trim(value), nil
			}
			// wasn't 'key=value' after all
			// TODO: use 'index' attribute from <name>
			return "", "", nil
		}
	}
}

func tag(s string) xml.Name {
	return xml.Name{Space: "", Local: s}
}

var trim = strings.TrimSpace
