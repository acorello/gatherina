package mirdin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"

	"dev.acorello.it/go/gatherina/jstree"
	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
	"github.com/acorello/must"
	"github.com/dop251/goja/ast"
)

type Newsletters struct {
	ArchEngineer map[string]string
}

func TestFetchArchEngineerNewsletter(t *testing.T) {
	var conf Newsletters
	_, err := toml.DecodeFile("newsletters.toml", &conf)
	if err != nil {
		log.Fatal(err)
	}
	res, err := http.Get(conf.ArchEngineer["jshook"])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	program := must.Get(jstree.Parse(res.Body))

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
