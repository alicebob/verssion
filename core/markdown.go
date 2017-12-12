// only knows about "[link](https://..)" markdown markup

package core

import (
	"fmt"
	"html/template"
	"strings"
)

type link struct {
	title, href, raw string
}

// parseMarkdown returns a list of strings or link{}s.
func parseMarkdown(s string) []interface{} {
	var (
		res   []interface{}
		start = 0
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '[':
			l, n := readLink(s[i:])
			if l != nil {
				if i > start {
					res = append(res, s[start:i])
				}
				start = i + n
				i--
				res = append(res, *l)
			}
			i += n
		}
	}
	if len(s)-start > 0 {
		res = append(res, s[start:])
	}
	return res
}

// readLink reads a "[...](...)" link at the start of s.
// Returns that link (if any), and how many bytes are used.
func readLink(s string) (*link, int) {
	title, n := readPair(s, '[', ']')
	if title == "" {
		return nil, 0
	}
	href, m := readPair(s[n:], '(', ')')
	if href == "" {
		return nil, 0
	}
	return &link{title: title, href: href, raw: s[:n+m]}, n + m
}

func readPair(s string, open, close byte) (string, int) {
	if len(s) == 0 || s[0] != open {
		return "", 0
	}
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++
		case open:
			return "", 0
		case close:
			res := strings.Replace(s[1:i], `\`, ``, -1)
			return res, i + 1
		}
	}
	return "", 0
}

// markdown to html (properly escaped)
func HtmlMarkdown(src string) template.HTML {
	var res string
	for _, t := range parseMarkdown(src) {
		switch s := t.(type) {
		case string:
			res += template.HTMLEscapeString(s)
		case link:
			if !safeURL(s.href) {
				res += template.HTMLEscapeString(s.raw)
				break
			}
			res += fmt.Sprintf(`<a href="%s">%s</a>`,
				template.HTMLEscapeString(s.href),
				template.HTMLEscapeString(s.title),
			)
		default:
			// panic?
		}
	}
	return template.HTML(res)
}

func safeURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// markdown to plain text
func TextMarkdown(src string) string {
	var res string
	for _, t := range parseMarkdown(src) {
		switch s := t.(type) {
		case string:
			res += s
		case link:
			res += s.title
		default:
			// panic?
		}
	}
	return res
}
