package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/acorello/must"
	"github.com/carlmjohnson/requests"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//go:embed query.json
var queryTempl string

func queryForPage(page int) string {
	return must.Get(sjson.Set(queryTempl, "searches[0].page", page))
}

const k = "IpJoOPZUczKNxR54gCnU8sjVNGCyXj21"

func main() {
	url := fmt.Sprintf("https://94eo8lmsqa0nd3j5p.a1.typesense.net/multi_search?x-typesense-api-key=%s", k)

	output := outputFile(cli.outputPath)
	defer output.Close()

	var page = 1
	for {
		fmt.Println("page", page)
		payload := queryForPage(page)

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
			log.Fatalf("error fetching page %d: %v", page, err)
		}

		hits := gjson.GetBytes(responseBody.Bytes(), "results.#.hits|@flatten")
		if !hits.IsArray() {
			log.Fatalf("error extracting data from page %d: unexpected type %s", page, hits.Type)
		}
		if len(hits.Array()) == 0 {
			fmt.Println("finished")
			break
		}
		for _, el := range hits.Array() {
			_, err = fmt.Fprintln(output, el.Get("document").Raw)
			if err != nil {
				log.Fatalf("error writing data of page %d: %v", page, err)
			}
		}
		page++
		if page > cli.maxPages {
			break
		}
		responseBody.Reset()
		time.Sleep(1 * time.Second)
	}
}

// outputFile creates and returns an *os.File for writing output data, based on the provided cliOutputPath.
// If cliOutputPath is empty, a temporary file is created. If it is a directory, a temp file is created within it.
// If it is a file, it ensures no overwrite occurs. It also handles creating missing parent directories.
func outputFile(cliOutputPath string) *os.File {
	var outputPath string
	if cliOutputPath == "" {
		tempFile := must.Get(os.CreateTemp("", "bcorps_results_*.jsonl"))
		tempFile.Close() // will be reopened again when writing
		outputPath = tempFile.Name()
	} else if s, err := os.Stat(cliOutputPath); err == nil {
		if s.IsDir() {
			outputPath = must.Get(os.MkdirTemp(cliOutputPath, "results_*.jsonl"))
		} else {
			log.Fatalf("file already exists: %s", cliOutputPath) // no auto-override
		}
	} else {
		outputPath = cliOutputPath
		err := os.MkdirAll(filepath.Base(cliOutputPath), 0750)
		if err != nil {
			log.Fatalf("error creating parent folders of %s: %v", cliOutputPath, err)
		}
	}

	fmt.Println("output path:", outputPath)
	output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	return output
}
