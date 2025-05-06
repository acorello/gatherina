package main

import (
	"io"
	"log"
	"net/url"
	"os"
)

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
