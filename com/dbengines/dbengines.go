// gathers DBMS details from db-engines.com
package dbengines

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Database struct {
	// CORE DETAILS

	// Expected format: "https://db-engines.com/$LANG/system/$DBNAME"
	DetailsPage url.URL `json:"-"`
	Name        string

	// ADDITIONAL DETAILS

	CloudOnly     string
	ImplementedIn string
	License       string
	PrimaryModel  string
}

func (my *Database) readDetails(d map[string]string) {
	my.CloudOnly = d["CLOUD-BASED ONLY"]
	my.ImplementedIn = d["IMPLEMENTATION LANGUAGE"]
	my.License = d["LICENSE"]
	my.PrimaryModel = d["PRIMARY DATABASE MODEL"]
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
		addDetails(&db)
		dbJSON, _ := json.Marshal(&db)
		fmt.Println(string(dbJSON))
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
	db.DetailsPage = *url
	return db, nil
}

func addDetails(db *Database) {
	details, err := getDetailsMap(db.DetailsPage)
	if err != nil {
		log.Println("failed to get details", err)
		return
	}
	db.readDetails(details)
}

func getDetailsMap(url url.URL) (details map[string]string, err error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	detailsDocument, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	attributeCells := detailsDocument.Find("TABLE.tools TBODY").Find("TD.attribute:not(.external_att)")
	details = make(map[string]string, attributeCells.Length())
	attributeCells.Each(func(i int, attributeTD *goquery.Selection) {
		attributeTD.Find("SPAN.info").Remove()
		name := strings.ToUpper(strings.TrimSpace(attributeTD.Text()))
		valueTD := attributeTD.Next()
		valueTD.Find("SPAN.info").Remove()
		value := strings.TrimSpace(valueTD.Text())
		if name == "" || value == "" {
			log.Printf("empty attribute %q or value %q\n", name, value)
			return
		}
		details[name] = value
	})
	return details, nil
}
