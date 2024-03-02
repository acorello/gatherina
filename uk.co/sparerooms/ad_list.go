package main

import (
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/url"
	"os"
	"text/template"
)

const adTempl = `
     Title {{.Title}}
Listing-Id {{.ListingId}}
    PostCo {{.PostCo}}
      Href {{.Href}}
`

var adTemplate = Must(template.New("ad").Parse(adTempl))

func PrintAdList(r io.Reader) {
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}
	articles := d.Find("li[class='listing-result']")
	var ads = make([]Ad, articles.Size())
	articles.Each(func(i int, article *goquery.Selection) {
		ads[i] = queryAd(article)
	})
	for i := range ads {
		adTemplate.Execute(os.Stdout, ads[i])
	}
}

type Ad struct {
	Title     string
	ListingId string
	PostCo    string
	Href      string
}

var firstFigureAnchor = goquery.Single("figure > a")

func queryAd(article *goquery.Selection) Ad {
	href, _ := url.Parse(article.FindMatcher(firstFigureAnchor).AttrOr("href", ""))
	q := href.Query()
	q.Del("search_results")
	href.RawQuery = q.Encode()
	ad := Ad{
		Title:     article.AttrOr("data-listing-title", "N/A"),
		ListingId: article.AttrOr("data-listing-id", "N/A"),
		PostCo:    article.AttrOr("data-listing-postcode", "N/A"),
		Href:      href.String(),
	}
	return ad
}
