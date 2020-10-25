package linkParser

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

//will take in an HTML document
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := nodeLinks(doc)

	return nodes, nil
}

func buildLink(node *html.Node) Link {
	var href, text string

	for _, attr := range node.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}

	text = findText(node)

	return Link{href, text}
}

func findText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	innerText := ""

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		innerText = fmt.Sprintf("%s %s", innerText, findText(child))
	}

	return strings.Join(strings.Fields(innerText), " ")
}

func nodeLinks(node *html.Node) []Link {
	if node.Type == html.ElementNode && node.Data == "a" {
		return []Link{buildLink(node)}
	}

	var links []Link

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		links = append(links, nodeLinks(child)...)
	}

	return links
}
