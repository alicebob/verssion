package w

import (
	"bytes"
)

type Page struct {
	StableVersion string
}

// GetPage downloads and parses given wikipage
func GetPage(rev int) (Page, error) {
	ptXML, err := GetParseTree(rev)
	if err != nil {
		return Page{}, nil
	}

	pt, err := Parse(bytes.NewBufferString(ptXML))
	if err != nil {
		return Page{}, err
	}

	return Page{
		StableVersion: StableVersion(pt),
	}, nil
}

func StableVersion(ts []Template) string {
	for _, t := range ts {
		switch t.Title {
		case "Infobox software":
			return t.NamedParts["latest release version"]
		case "LSR":
			return t.NamedParts["latest_release_version"]
		}
	}
	return ""
}
