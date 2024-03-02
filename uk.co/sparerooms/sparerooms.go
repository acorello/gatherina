package main

import (
	"io"
	"log"
	"net/url"
	"os"
)

func main() {
	searchURL := *Must(url.Parse("/Users/am/Projects/gatherina/uk.co/sparerooms/samples/search@01.html"))
	processURL(searchURL, PrintAdList)
}

func processURL(u url.URL, consumer func(io.Reader)) {
	if !(u.Scheme == "file" || u.Scheme == "") {
		log.Fatal("Only file URLs currently supported")
	}
	f, err := os.Open(u.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	consumer(f)
}

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return t
}
