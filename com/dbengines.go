// gathers DBMS details from db-engines.com
package dbengines

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Database struct {
	// Expected format: "https://db-engines.com/$LANG/system/$DBNAME"
	URL  url.URL
	Name string
}

const (
	HostName    = "db-engines.com"
	systemsPage = "https://" + HostName + "/en/systems"
)

func DatabaseDetail() {
	systemsDocument := getSystemsDocument()

	linksSelection := systemsDocument.Find("table.list a")
	databases := make([]Database, 0, linksSelection.Length())
	linksSelection.Each(func(_ int, dbAnchor *goquery.Selection) {
		db, err := dbLink(dbAnchor)
		if err != nil {
			log.Println("invalid item:", err)
			return
		}
		databases = append(databases, db)
	})
	for _, db := range databases {
		fmt.Println(db.Name, db.URL)
	}
}

func getSystemsDocument() *goquery.Document {
	resp, err := http.Get(systemsPage)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	systemsPage, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return systemsPage
}

func dbLink(anchor *goquery.Selection) (db Database, err error) {
	anchor.Find("span").Remove()
	db.Name = strings.TrimSpace(anchor.Text())
	if db.Name == "" {
		return db, fmt.Errorf("missing 'name'")
	}
	href, _ := anchor.Attr("href")
	href = strings.TrimSpace(href)
	if href == "" {
		return db, fmt.Errorf("missing 'href'")
	}
	url, urlErr := url.Parse(href)
	if urlErr != nil {
		return db, fmt.Errorf("invalid 'href' for %q: %s", db.Name, urlErr)
	}
	db.URL = *url
	return db, nil
}
