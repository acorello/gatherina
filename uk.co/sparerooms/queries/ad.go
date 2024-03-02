package queries

import (
	"dev.acorello.it/go/gatherina/must"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"regexp"
	"strconv"
)

type AdDetails struct {
	ListingId string
	Location
}

type Location struct {
	Latitude, Longitude float64
}

var locationRegEx = must.Must(regexp.Compile(`location\s*:\s*\{\s*latitude:\s*"([^"]+)"\s*,\s*longitude\s*:\s*"([^"]+)",}`))

func Ad(r io.Reader) (ad AdDetails, err error) {
	d, gqErr := goquery.NewDocumentFromReader(r)
	if gqErr != nil {
		return ad, gqErr
	}
	if lat, lon, err := queryLocation(d); err != nil {
		return ad, err
	} else {
		ad.Latitude = lat
		ad.Longitude = lon
	}
	if id, err := queryAdId(d); err != nil {
		return ad, err
	} else {
		ad.ListingId = id
	}
	return ad, nil
}
func queryAdId(d *goquery.Document) (id string, err error) {
	err = fmt.Errorf("ad-identifier not found")
	shareDiv := d.Find("div#share")
	if shareDiv == nil {
		return id, err
	}
	id, found := shareDiv.Attr("data-advert-id")
	if !found {
		return id, err
	}
	return id, nil
}

func queryLocation(d *goquery.Document) (lat float64, lon float64, err error) {
	err = fmt.Errorf("location not found")

	scripts := d.Find("script")
	parseFloat := func(s string) float64 {
		return must.Must(strconv.ParseFloat(s, 64))
	}
	scripts.EachWithBreak(func(i int, script *goquery.Selection) (cont bool) {
		js := script.Text()
		if match := locationRegEx.FindStringSubmatch(js); match != nil {
			lat = parseFloat(match[1])
			lon = parseFloat(match[2])
			err = nil
			return false
		}
		return true
	})
	return lat, lon, err
}
