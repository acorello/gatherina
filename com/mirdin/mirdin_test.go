package mirdin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"dev.acorello.it/go/gatherina/jstree"
	"github.com/PuerkitoBio/goquery"
	"github.com/dop251/goja/ast"
)

func TestFetchArchEngineerNewsletter(t *testing.T) {
	const URL = "https://mirdin.us16.list-manage.com/generate-js/?u=8b565c97b838125f69e75fb7f&amp;fid=313471&amp;show=100000000"
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	program := jstree.MustParse(res.Body)

	var html string
	jstree.Walk(program, func(n ast.Node, _ int) bool {
		switch node := n.(type) {
		case *ast.StringLiteral:
			json.Unmarshal([]byte(node.Literal), &html)
			return false
		default:
			return true
		}
	})
	mailchimpNewsLinks, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	links := mailchimpNewsLinks.Find("a")
	if !(links.Length() > 0) {
		t.Fatal("no links found")
	}
	links.Each(func(i int, s *goquery.Selection) {
		href, found := s.Attr("href")
		if found {
			fmt.Println(s.Text(), "\n\t", href)
		}
	})
}
