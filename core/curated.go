package core

import (
	"fmt"
	"strings"
	"time"
)

type Curated struct {
	CustomTitle string
	Created     time.Time
	LastUsed    time.Time
	LastUpdated time.Time
	Pages       []string
}

func (c *Curated) Title() string {
	if t := c.CustomTitle; t != "" {
		return t
	}
	return c.DefaultTitle()
}

func (c *Curated) DefaultTitle() string {
	var (
		p        = c.Pages
		ellipses = 0
		maxEls   = 4
	)
	if len(p) == 0 {
		return "[untitled feed]"
	}
	if len(p) > maxEls {
		ellipses = len(p) - maxEls
		p = p[:maxEls]
	}
	t := strings.Join(Titles(p), ", ")
	if ellipses > 0 {
		t += fmt.Sprintf(", ... (%d more)", ellipses)
	}
	return t
}
