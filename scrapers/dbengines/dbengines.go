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
	DetailsPage url.URL `json:"-"` // format: "https://db-engines.com/$LANG/system/$DBNAME"
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
	knownDBPage := fetchKnownDBPage()
	dbList := knownDBPage.Find("table.list a")
	databases := listedDatabaseWithDetailsURL(dbList)
	for _, db := range databases {
		addDetails(&db)
		dbJSON, _ := json.Marshal(&db)
		fmt.Println(string(dbJSON))
	}
}

func listedDatabaseWithDetailsURL(linksSelection *goquery.Selection) []Database {
	databases := make([]Database, linksSelection.Length())
	linksSelection.Each(func(i int, dbAnchor *goquery.Selection) {
		var db Database
		err := setNameAndURL(&db, dbAnchor)
		if err != nil {
			log.Printf("invalid item %d: %v", i, err)
			return
		}
		databases[i] = db
	})
	return databases
}

func fetchKnownDBPage() *goquery.Document {
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

func setNameAndURL(db *Database, anchor *goquery.Selection) (err error) {
	anchor.Find("span").Remove()
	db.Name = strings.TrimSpace(anchor.Text())
	if db.Name == "" {
		return fmt.Errorf("missing 'name'")
	}
	href, found := anchor.Attr("href")
	if !found {
		return fmt.Errorf(`missing "href" attribute`)
	}
	aURL, err := parseURL(href)
	if err != nil {
		return fmt.Errorf(`failed to parse "href": %w`, err)
	}
	db.DetailsPage = *aURL
	return nil
}

func parseURL(href string) (dbURL *url.URL, err error) {
	href = strings.TrimSpace(href)
	if href == "" {
		return dbURL, fmt.Errorf("blank")
	}
	dbURL, err = url.Parse(href)
	if err != nil {
		return dbURL, err
	}
	return dbURL, nil
}

func addDetails(db *Database) {
	details, err := fetchDetailsMap(db.DetailsPage)
	if err != nil {
		log.Println("failed to get details", err)
		return
	}
	db.readDetails(details)
}

func fetchDetailsMap(url url.URL) (details map[string]string, err error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	attributeNodes := doc.Find("TABLE.tools TBODY").Find("TD.attribute:not(.external_att)")
	details = make(map[string]string, attributeNodes.Length())
	attributeNodes.Each(func(i int, attributeNode *goquery.Selection) {
		name, value := dbAttribute(attributeNode)
		if name == "" || value == "" {
			log.Printf("[fetchDetailsMap] empty attribute %q or value %q for item %d", name, value, i)
			return
		}
		details[name] = value
	})
	return details, nil
}

func dbAttribute(attributeNode *goquery.Selection) (name string, value string) {
	attributeNode.Find("SPAN.info").Remove()
	name = strings.ToUpper(strings.TrimSpace(attributeNode.Text()))
	valueTD := attributeNode.Next()
	valueTD.Find("SPAN.info").Remove()
	value = strings.TrimSpace(valueTD.Text())
	return name, value
}
