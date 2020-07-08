package main

import (
	"fmt"

	"github.com/antchfx/htmlquery"
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

	// analysis page content
	doc, err := htmlquery.Parse(reader)

	pageTitleNode := htmlquery.FindOne(doc, `/html/head/title`)
	pageTitle := htmlquery.InnerText(pageTitleNode)

	if t.gc.firstParse {
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
		t.gc.firstParse = false
	}

	rules := ctx.config.Rules

	for _, r := range rules {
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

					ext, err := getExtension(urlString)
					if err != nil {
						ctx.log.Errorf("Cannot parse url in download file task, err: %s", err.Error())
						return err
					}
					nc := &namingContext{
						pageTitle: pageTitle,
						imgAlt:    htmlquery.SelectAttr(n, "alt"),
						text:      htmlquery.InnerText(n),
						ext:       ext,
						i:         t.gc.i,
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
