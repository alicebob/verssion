package web

import (
	"net/url"
)

func adhocURL(base string, pages []string) string {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	u.Path += "/adhoc/atom.xml"
	u.RawQuery = url.Values{
		"p": pages,
	}.Encode()
	return u.String()
}
