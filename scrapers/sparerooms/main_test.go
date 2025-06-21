package main

import (
	"net/url"
	"testing"

	"github.com/acorello/must"
)

func TestProcessingSample(t *testing.T) {
	searchURL := *must.Get(url.Parse("testdata/samples/search@01.html"))
	processURL(searchURL, PrintAdList)
}
