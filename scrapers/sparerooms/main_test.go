package main

import (
	"net/url"
	"testing"

	"github.com/acorello/must"
	"github.com/matryer/is"
)

func TestProcessingSample(t *testing.T) {
	searchURL := *must.Get(url.Parse("testdata/samples/search@01.html"))
	processURL(searchURL, PrintAdList)
}

func TestIs(t *testing.T) {
	st := is.New(t)
	st.Equal(1, 2) // expected equality
	t.Error("not equal")
}
