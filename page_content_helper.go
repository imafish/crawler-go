package main

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/antchfx/htmlquery"
	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/net/html"
)

// GetPageTitle returns page title
func GetPageTitle(doc *html.Node, converter *iconv.Converter) string {
	pageTitle := ""
	pageTitleNode := htmlquery.FindOne(doc, `/html/head/title`)
	if pageTitleNode != nil {
		pageTitle = htmlquery.InnerText(pageTitleNode)
	}

	pageTitle = convertEncoding(pageTitle, converter)

	return pageTitle
}

// GetPageEncoding returns charset property in html meta field, if it exists.
// It returns empty string if charset property doesnot exist.
func GetPageEncoding(doc *html.Node) string {
	encoding := ""
	meta := htmlquery.FindOne(doc, `/html/head/meta`)
	if meta != nil {

		encodingString := htmlquery.SelectAttr(meta, "content")
		reg := regexp.MustCompile(`charset=(\w+)`)
		matches := reg.FindStringSubmatch(encodingString)
		if len(matches) == 2 {
			encoding = matches[1]
		}
	}

	return encoding
}

// GetMatchingNodes returns matching html nodes.
func GetMatchingNodes(top *html.Node, t target) ([]*html.Node, error) {
	var matchingNodes []*html.Node

	switch {
	case t.Xpath != "":
		matchingNodes = htmlquery.Find(top, t.Xpath)
	default:
		return nil, fmt.Errorf("Unsupported target type: %#v", t)
	}

	return matchingNodes, nil
}

// CollectText returns all inner texts of one node (and its child nodes)
func CollectText(node *html.Node, converter *iconv.Converter) string {
	buf := &bytes.Buffer{}

	CollectTextRecursive(node, buf)
	text := buf.String()
	text = convertEncoding(text, converter)

	return text
}

// CollectTextRecursive recursively collects all text from node.
func CollectTextRecursive(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		CollectTextRecursive(c, buf)
	}
}

func convertEncoding(in string, converter *iconv.Converter) string {
	str := in
	if converter != nil {
		var err error
		str, err = converter.ConvertString(in)
		if err != nil {
			str = in
		}
	}

	return str
}
