package ft

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func TestHTMLExtractor(t *testing.T) {
	// https://www.ft.com/ft1000-2024
	// #article-body > div.o-table-container.n-content-layout.o-table-container--expanded > div > div.o-table-scroll-wrapper > table
	// //*[@id="article-body"]/div[1]/div/div[1]/table
	/*
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
	input, err := os.Open(`testdata/FT1000-ranking.fmt.html`)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		panic(err)
	}
	type Record struct {
		Rank int
		Name,
		Country,
		Sector,
		AbsoluteGrowthRatePct,
		CompoundAnnualGrowthRatePct,
		Revenue2022Eur,
		Revenue2019Eur,
		Employees2022,
		Employees2019,
		FoundingYear string
		Url string
	}

	for _, tr := range doc.Find(`table tbody tr`).EachIter() {
		var values Record
		cells := tr.Find(`td`)
		rank := parseInt(text(cells.Get(0)))
		values.Rank = rank
		values.Name = strings.TrimSpace(allText(cells.Get(1)))
		values.Url = cells.Find(`a`).AttrOr(`href`, "")
		values.Country = text(cells.Get(2))
		values.Sector = text(cells.Get(3))
		values.AbsoluteGrowthRatePct = text(cells.Get(4))
		values.CompoundAnnualGrowthRatePct = text(cells.Get(5))
		values.Revenue2022Eur = text(cells.Get(6))
		values.Revenue2019Eur = text(cells.Get(7))
		values.Employees2022 = text(cells.Get(8))
		values.Employees2019 = text(cells.Get(9))
		values.FoundingYear = text(cells.Get(10))
		fmt.Printf("%#v\n", values)
	}
}

func parseInt(s string) int {
	rank, err := strconv.Atoi(s)
	if err != nil {
		fmt.Errorf("failed to parse as integer, returning default; got: %q", s)
		return 0
	}
	return rank
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
