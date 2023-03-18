// gathers DBMS details from db-engines.com
package dbengines

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Database struct {
	URL  string
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
		db := dbLink(dbAnchor)
		if db.Name == "" || db.URL == "" {
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

func dbLink(anchor *goquery.Selection) Database {
	anchor.Find("span").Remove()
	href, _ := anchor.Attr("href")
	href = strings.TrimSpace(href)
	name := strings.TrimSpace(anchor.Text())
	return Database{
		URL:  href,
		Name: name,
	}
}
