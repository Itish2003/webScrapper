package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type WikiInfo struct {
	Headings string
	Content  string
}

func main() {

	wiki := WikiInfo{}
	scrapeUrl := "https://en.wikipedia.org/wiki/Wiki"

	c := colly.NewCollector(colly.AllowedDomains("www.en.wikipedia.org", "en.wikipedia.org"))

	file, err := os.Create("scraper.txt")
	if err != nil {
		fmt.Println("File not created.")
	}

	writer := bufio.NewWriter(file)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,q=0.9")
		writer.WriteString(r.URL.String())
	})

	c.OnHTML("p", func(h *colly.HTMLElement) {
		wiki.Headings = h.Text
		writer.WriteString(h.Text)
	})

	c.OnError(func(r *colly.Response, e error) {
		writer.WriteString(e.Error())
	})

	c.OnScraped(func(c *colly.Response) {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", " ")
		enc.Encode(wiki)
	})

	writer.Flush()
	c.Visit(scrapeUrl)
	defer file.Close()
}
