package main

import (
	"io"
	"log"
	"net/url"
	"os"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/acorello/must"
)

const adTempl = `
     Title {{.Title}}
Listing-Id {{.ListingId}}
 Post-Code {{.PostCode}}
       URL {{.HRef}}
`

var adTemplate = must.Get(template.New("ad").Parse(adTempl))

func PrintAdList(r io.Reader) {
	ads, err := AdList(r)
	if err != nil {
		log.Fatal(err)
	}
	for i := range ads {
		adTemplate.Execute(os.Stdout, ads[i])
	}
}

func AdList(r io.Reader) (ads []AdSummary, err error) {
	d, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	articles := d.Find("li[class='listing-result']")
	ads = make([]AdSummary, articles.Size())
	articles.Each(func(i int, article *goquery.Selection) {
		ads[i] = queryAdSummary(article)
	})
	return ads, nil
}

type AdSummary struct {
	Title     string
	ListingId string
	PostCode  string
	HRef      string
}

var firstFigureAnchor = goquery.Single("figure > a")

func queryAdSummary(article *goquery.Selection) AdSummary {
	href, _ := url.Parse(article.FindMatcher(firstFigureAnchor).AttrOr("href", ""))
	q := href.Query()
	q.Del("search_results")
	href.RawQuery = q.Encode()
	ad := AdSummary{
		Title:     article.AttrOr("data-listing-title", "N/A"),
		ListingId: article.AttrOr("data-listing-id", "N/A"),
		PostCode:  article.AttrOr("data-listing-postcode", "N/A"),
		HRef:      href.String(),
	}
	return ad
}
