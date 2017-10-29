package main

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
`)
	want := []string{
		"Foo1", "Foo2", "http://En.wiKIPedia.oRg/wiKI/Foo3", "Foo4", "Foo5", "wiki/Foo6",
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have %v, want %v", have, want)
	}
}
