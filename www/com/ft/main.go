package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func main() {
	input, err := os.Open(`testdata/FT1000-ranking.fmt.html`)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		panic(err)
	}
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	var csvHeader []string
	for _, m := range mappings {
		csvHeader = append(csvHeader, m.Name)
	}
	err = csvWriter.Write(csvHeader)
	if err != nil {
		log.Fatal(err)
	}
	for _, tr := range doc.Find(`table tbody tr`).EachIter() {
		cells := tr.Find(`td`)

		var csvRow []string
		for _, m := range mappings {
			csvRow = append(csvRow, m.GetFrom(cells))
		}

		err = csvWriter.Write(csvRow)
		csvRow = csvRow[:0]
		if err != nil {
			log.Fatal(err)
		}
	}
}

// https://www.ft.com/ft1000-2024
// #article-body > div.o-table-container.n-content-layout.o-table-container--expanded > div > div.o-table-scroll-wrapper > table
// //*[@id="article-body"]/div[1]/div/div[1]/table
/*
	// THEAD
	<th data-o-table-data-type="text" scope="col" role="columnheader">
	  <strong>Rank</strong>
	</th>
	<th data-o-table-data-type="text" scope="col" role="columnheader">
	  <strong>Name</strong>
	</th>
	<th data-o-table-data-type="text" scope="col" role="columnheader">
	  <strong>Country</strong>
	</th>
	<th data-o-table-data-type="text" scope="col" role="columnheader">
	  <strong>Sector</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Absolute Growth Rate (%)</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Compound Annual Growth Rate (%)</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Revenue 2022 (€)</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Revenue 2019 (€)</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Number of employees 2022</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Number of employees 2019</strong>
	</th>
	<th data-o-table-data-type="number" scope="col" role="columnheader">
	  <strong>Founding Year</strong>
	</th>
	// TBODY
	<tr aria-hidden="false">
	  <td class="">1000</td>
	  <td class="">
		<a href="https://www.raketspaservices.com" data-trackable="link">Rak Et Spa Services</a>
	  </td>
	  <td class="">France</td>
	  <td class=""><![CDATA[Food & Beverages]]></td>
	  <td class="o-table__cell--numeric">156.5%</td>
	  <td class="o-table__cell--numeric">36.9%</td>
	  <td class="o-table__cell--numeric">16,639,248</td>
	  <td class="o-table__cell--numeric">6,487,824</td>
	  <td class="o-table__cell--numeric">22</td>
	  <td class="o-table__cell--numeric">10</td>
	  <td class="o-table__cell--numeric">2008</td>
	</tr>
*/
var mappings = []ColMapping{
	{Name: "Rank", TDIdx: 0},
	{Name: "Name", TDIdx: 1, Extractor: func(idx int, cells *goquery.Selection) string {
		return strings.TrimSpace(allText(cells.Get(1)))
	}},
	{Name: "Url", TDIdx: 1, Extractor: func(idx int, cells *goquery.Selection) string {
		td := &goquery.Selection{Nodes: []*html.Node{cells.Get(1)}}
		return td.Find(`a`).AttrOr(`href`, "")
	}},
	{Name: "Country", TDIdx: 2},
	{Name: "Sector", TDIdx: 3},
	{Name: "AbsoluteGrowthRatePct", TDIdx: 4},
	{Name: "CompoundAnnualGrowthRatePct", TDIdx: 5},
	{Name: "Revenue2022Eur", TDIdx: 6, Extractor: func(idx int, cells *goquery.Selection) string {
		return remove(text(cells.Get(6)), ",")
	}},
	{Name: "Revenue2019Eur", TDIdx: 7, Extractor: func(idx int, cells *goquery.Selection) string {
		return remove(text(cells.Get(7)), ",")
	}},
	{Name: "Employees2022", TDIdx: 8},
	{Name: "Employees2019", TDIdx: 9},
	{Name: "FoundingYear", TDIdx: 10},
}

type ColMapping struct {
	Name      string
	TDIdx     int
	Extractor func(idx int, row *goquery.Selection) string
}

func (m ColMapping) GetFrom(cells *goquery.Selection) string {
	if m.Extractor == nil {
		return text(cells.Get(m.TDIdx))
	}
	return m.Extractor(m.TDIdx, cells)
}

func remove(s string, unwanted ...string) string {
	for _, u := range unwanted {
		s = strings.ReplaceAll(s, u, "")
	}
	return s
}

func text(n *html.Node) string {
	var builder strings.Builder
	if n.Type == html.TextNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			// Keep newlines and spaces, like jQuery
			builder.WriteString(c.Data)
		}
	}
	return builder.String()
}

func allText(n *html.Node) string {
	var builder strings.Builder
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			builder.WriteString(n.Data)
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for n := range n.ChildNodes() {
		f(n)
	}
	return builder.String()
}
