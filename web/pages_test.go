package web

import (
	"reflect"
	"testing"
)

func TestToPages(t *testing.T) {
	haveOK, haveErr := toPages(`
http://en.wikipedia.org/wiki/Foo1
http://En.wiKIPedia.oRg/wiki/Foo2
http://En.wiKIPedia.oRg/wiKI/Foo3
https://en.wikipedia.org/wiki/Foo4
Foo5

wiki/Foo6
Foo 7
`)
	wantOK := []string{
		"Foo1", "Foo2", "http://En.wiKIPedia.oRg/wiKI/Foo3", "Foo4", "Foo5", "wiki/Foo6",
	}
	wantErr := []string{
		`invalid page: "Foo 7"`,
	}
	if !reflect.DeepEqual(haveOK, wantOK) {
		t.Errorf("have %v, want %v", haveOK, wantOK)
	}
	var have []string
	for _, e := range haveErr {
		have = append(have, e.Error())
	}
	if !reflect.DeepEqual(have, wantErr) {
		t.Errorf("have %v, want %v", have, wantErr)
	}
}
