package web

import (
	"reflect"
	"testing"
)

func TestToPages(t *testing.T) {
	have := toPages(`
http://en.wikipedia.org/wiki/Foo1
http://En.wiKIPedia.oRg/wiki/Foo2
http://En.wiKIPedia.oRg/wiKI/Foo3
https://en.wikipedia.org/wiki/Foo4
Foo5

wiki/Foo6
Foo 7
/wiki/Foo8
`)
	want := []string{
		"Foo1",
		"Foo2",
		"/wiKI/Foo3",
		"Foo4",
		"Foo5",
		"wiki/Foo6",
		"Foo_7",
		"Foo8",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}
