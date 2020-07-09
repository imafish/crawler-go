package main

import (
	"fmt"
	"regexp"

	"github.com/antchfx/htmlquery"
	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/net/html"
)

// PageParseTask holds information to parse a page
type PageParseTask struct {
	url   string
	final bool

	gc groupContext
}

// URL returns the task's corresponding url
func (t PageParseTask) URL() string {
	return t.url
}

// Execute grab and analyse the page content, generates other related tasks.
func (t PageParseTask) Execute(ctx *Context) error {
	reader, err := grabPageReader(t.url)
	if err != nil {
		ctx.log.Errorf("failed to grab page from url %s, err: %s", t.url, err.Error())
		return err
	}
	doc, err := htmlquery.Parse(reader)
	if err != nil {
		ctx.log.Errorf("failed to parse page content. err: %s", err.Error())
		return err
	}

	if t.gc.firstParse {
		// get page encoding and initialize converter
		var converter *iconv.Converter
		meta := htmlquery.FindOne(doc, `/html/head/meta`)
		if meta != nil {
			encoding := ""
			encodingString := htmlquery.SelectAttr(meta, "content")
			reg := regexp.MustCompile(`charset=(\w+)`)
			matches := reg.FindStringSubmatch(encodingString)
			if len(matches) == 2 {
				encoding = matches[1]
			}

			if encoding != "" && encoding != "utf-8" {
				converter, err = iconv.NewConverter(encoding, "utf-8")
				if err != nil {
					ctx.log.Errorf("failed to create encoding converter from %s. err: %s", encoding, err.Error())
					return err
				}
			}
		}
		t.gc.converter = converter

		// dir formatting for group context
		pageTitle := getPageTitle(doc)
		pageTitle = cc(pageTitle, converter)
		nc := &namingContext{
			i:         t.gc.i,
			pageTitle: pageTitle,
		}
		dir, err := formatString(t.gc.dir, nc, ctx)
		if err != nil {
			ctx.log.Errorf("failed to parse dir string %s, err: %s", t.gc.dir, err.Error())
			return err
		}
		t.gc.dir = dir

		ctx.log.Debugf("generated group context: %#v", t.gc)
		t.gc.firstParse = false
	}

	rules := ctx.config.Rules

	for _, r := range rules {
		// skip rule if only process link action and start group action present, and 'Final' is true,
		if r.Action.DownloadFile == nil && t.final {
			continue
		}

		for _, tgt := range r.Targets {
			matchingNodes, err := getMatchingNodes(doc, tgt)
			if err != nil {
				ctx.log.Errorf("Unable to get matching nodes. err: %s", err.Error())
				return err
			}
			for _, n := range matchingNodes {
				ctx.log.Debugf("Find a matching node: %v", n)
				if df := r.Action.DownloadFile; df != nil {
					urlString, err := getOneNodeInnerText(n, df.Target)
					if err != nil {
						ctx.log.Errorf("Failed to get inner text of sub-target, err: %s", err.Error())
						return err
					}
					urlString, err = convertToAbsoluteURL(t.url, urlString)
					if err != nil {
						ctx.log.Errorf("failed to convert url %s to absolute", urlString)
						return err
					}

					nc, err := createNamingContext(urlString, t, doc, n)
					if err != nil {
						ctx.log.Errorf("failed to create naming context. err: %s", err.Error())
						return err
					}

					task := DownloadFileTask{
						url:             urlString,
						dirPattern:      t.gc.dir,
						filenamePattern: df.FilenamePattern,
						nc:              nc,
					}
					ctx.addTask(task)
				}

				if r.Action.StartGroup != nil && !t.final {
					ctx.log.Error("StartGroup is not supported yet.")
					return fmt.Errorf("StartGroup is not supported yet")
				}

				if pl := r.Action.ProcessLink; pl != nil && !t.final {
					urlString, err := getOneNodeInnerText(n, pl.Target)
					if err != nil {
						ctx.log.Errorf("Failed to get inner text of sub-target, err: %s", err.Error())
						return err
					}
					if isInvalid(urlString) {
						continue
					}
					urlString, err = convertToAbsoluteURL(t.url, urlString)
					if err != nil {
						ctx.log.Errorf("Unable to convert to absolute path, err: %s", err.Error())
						return err
					}

					task := PageParseTask{
						url:   urlString,
						final: pl.Final,
						gc:    t.gc,
					}
					ctx.addTask(task)
				}
			}
		}
	}

	return nil
}

func getOneNodeInnerText(top *html.Node, t target) (string, error) {
	var subNode *html.Node
	switch {
	case t.Xpath != "":
		subNode = htmlquery.FindOne(top, t.Xpath)
	default:
		return "", fmt.Errorf("Unsupported target type: %#v", t)
	}
	if subNode == nil {
		return "", nil
	}
	urlString := htmlquery.InnerText(subNode)

	return urlString, nil
}

func getMatchingNodes(top *html.Node, t target) ([]*html.Node, error) {
	var matchingNodes []*html.Node

	switch {
	case t.Xpath != "":
		matchingNodes = htmlquery.Find(top, t.Xpath)
	default:
		return nil, fmt.Errorf("Unsupported target type: %#v", t)
	}

	return matchingNodes, nil
}

func isInvalid(url string) bool {
	invalid := false

	switch url {
	case "#":
		invalid = true
	case "":
		invalid = true
	default:
		invalid = false
	}

	return invalid
}

func getPageTitle(doc *html.Node) string {
	pageTitle := ""
	pageTitleNode := htmlquery.FindOne(doc, `/html/head/title`)
	if pageTitleNode != nil {
		pageTitle = htmlquery.InnerText(pageTitleNode)
	}

	return pageTitle
}

func cc(in string, converter *iconv.Converter) string {
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

func createNamingContext(url string, t PageParseTask, doc, n *html.Node) (*namingContext, error) {
	ext, err := getExtension(url)
	if err != nil {
		return nil, err
	}

	pageTitle := getPageTitle(doc)
	pageTitle = cc(pageTitle, t.gc.converter)

	imgAlt := htmlquery.SelectAttr(n, "alt")
	imgAlt = cc(pageTitle, t.gc.converter)

	text := htmlquery.InnerText(n)
	text = cc(pageTitle, t.gc.converter)

	nc := &namingContext{
		pageTitle: pageTitle,
		imgAlt:    imgAlt,
		text:      text,
		ext:       ext,
		i:         t.gc.i,
	}

	return nc, nil
}
