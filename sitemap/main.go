package main

import (
	"flag"
	"fmt"
	"github.com/w1kend/go/link"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	flagUrl := flag.String("url", "https://gophercises.com", "the url to build a sitemap")
	maxDepth := flag.Int("max-depth", 3, "max depth")
	flag.Parse()

	hrefs := bfs(*flagUrl, *maxDepth)

	for _, href := range hrefs {
		fmt.Println(href)
	}

	print(fmt.Sprintf("Founded %d links\n", len(hrefs)))
}

func bfs(siteUrl string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var queue map[string]struct{}

	nextQueue := map[string]struct{}{
		siteUrl: struct{}{},
	}

	for i := 0; i < maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{})

		for currUrl, _ := range queue {
			if _, ok := seen[currUrl]; ok {
				continue
			}

			seen[currUrl] = struct{}{}

			for _, nextUrl := range get(currUrl) {
				nextQueue[nextUrl] = struct{}{}
			}
		}
	}

	//get map keys
	result := make([]string, 0, len(seen))

	for seenUrl, _ := range seen {
		result = append(result, seenUrl)
	}

	return result
}

func get(siteUrl string) []string {
	response, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	requestUrl := response.Request.URL
	baseUrl := &url.URL{
		Scheme: requestUrl.Scheme,
		Host:   requestUrl.Host,
	}

	base := baseUrl.String()

	links, _ := link.Parse(response.Body)

	formattedHrefs := formatHrefs(links, base)

	return filter(base, formattedHrefs)
}

func formatHrefs(links []link.Link, baseUrl string) []string {
	var hrefs []string

	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, baseUrl+l.Href)
			break
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}

	return hrefs
}

func filter(base string, links []string) []string {
	var result []string

	for _, l := range links {
		if strings.HasPrefix(l, base) {
			result = append(result, l)
		}
	}

	return result
}
