package net_bcorporation

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/tidwall/gjson"
)

func TestFetchCompanies(t *testing.T) {
	apiKey := "IpJoOPZUczKNxR54gCnU8sjVNGCyXj21"
	url := fmt.Sprintf("https://94eo8lmsqa0nd3j5p.a1.typesense.net/multi_search?x-typesense-api-key=%s", apiKey)

	var page = 1
	output, err := os.OpenFile("results.jsonl", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer output.Close()
	for {
		fmt.Println("page", page)
		// JSON payload
		payload := fmt.Sprintf(`{
  "searches": [
    {
      "query_by": "name,description,websiteKeywords,countries,industry,sector,hqCountry,hqProvince,hqCity,hqPostalCode,provinces,cities,size,demographicsList",
      "exhaustive_search": true,
      "sort_by": "initialCertificationDateTimestamp:asc",
      "highlight_full_fields": "name,description,websiteKeywords,countries,industry,sector,hqCountry,hqProvince,hqCity,hqPostalCode,provinces,cities,size,demographicsList",
      "collection": "companies-production-en-us",
      "q": "*",
      "facet_by": "hqCountry",
      "max_facet_values": 500,
      "per_page": 50,
      "page": %d
    }
  ]
}`, page)

		var responseBody bytes.Buffer
		err := requests.New().Post().
			BaseURL(url).
			ContentType("text/plain").
			UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.4 Safari/605.1.15").
			Accept("application/json, text/plain, */*").
			Header("Sec-Fetch-Site", "cross-site").
			Header("Accept-Language", "en-GB,en;q=0.9").
			//Header("Accept-Encoding", "gzip, deflate, br").
			Header("Sec-Fetch-Mode", "cors").
			Header("Origin", "https://www.bcorporation.net").
			Header("Referer", "https://www.bcorporation.net/").
			Header("Sec-Fetch-Dest", "empty").
			Header("Priority", "u=3, i").
			BodyReader(strings.NewReader(payload)).
			ToBytesBuffer(&responseBody).
			Fetch(context.Background())
		if err != nil {
			t.Fatal(err)
		}

		hits := gjson.GetBytes(responseBody.Bytes(), "results.#.hits|@flatten")
		if !hits.IsArray() {
			t.Error("unexpected type", hits.Type)
		}
		if len(hits.Array()) == 0 {
			fmt.Println("finished")
			break
		}
		for _, el := range hits.Array() {
			_, err = fmt.Fprintln(output, el.Get("document").Raw)
			if err != nil {
				t.Fatal(err)
			}
		}
		page++
		//if page > 3 {
		//	break
		//}
		responseBody.Reset()
		time.Sleep(1 * time.Second)
	}
}
