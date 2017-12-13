package web

import (
	"github.com/alicebob/verssion/core"
)

func WikiSpider() core.Spider {
	return core.NewSpiderCache(core.NewWikipediaSpider(WikiURL))
}

func WikiURL(page string) string {
	return "https://en.wikipedia.org/wiki/" + page
}
