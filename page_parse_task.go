package main

import (
	"github.com/antchfx/htmlquery"
)

// PageParseTask holds information to parse a page
type PageParseTask struct {
	url   string
	final bool

	// for name formatting
	linkText string
}

// Execute grab and analyse the page content, generates other related tasks.
func (t PageParseTask) Execute(ctx *Context) error {
	goon := ctx.startExecuting(t.url)
	if !goon {
		ctx.log.Debugf("task already exists, skipping. task: %#v", t)
		return nil
	}

	reader, err := grabPageReader(t.url)
	if err != nil {
		ctx.log.Errorf("failed to grab page from url %s, err: %s", t.url, err.Error())
		return err
	}

	// analysis page content
	doc, err := htmlquery.Parse(reader)

	// add download file task to list:
	downloads := htmlquery.Find(doc, `//*[@id="bigimg"]`)
	for _, d := range downloads {
		urlString := htmlquery.SelectAttr(d, "src")
		url, err := convertToAbsoluteURL(t.url, urlString)
		if err != nil {
			ctx.log.Errorf("failed to convert url %s to absolute", urlString)
			return err
		}

		// getting context for name formatting
		pageTitleNode := htmlquery.FindOne(doc, `/html/head/title`)
		pageTitle := htmlquery.InnerText(pageTitleNode)
		imgAlt := htmlquery.SelectAttr(d, "alt")
		linkText := t.linkText

		task := DownloadFileTask{
			url:             url,
			dirPattern:      "img",
			filenamePattern: `{__i}{.ext}`,

			pageTitle: pageTitle,
			imgAlt:    imgAlt,
			linkText:  linkText,
		}

		ctx.addTask(task)
	}

	// add links to task list:
	if !t.final {
		links := htmlquery.Find(doc, "/html/body/div[19]/ul/li/a")
		for _, link := range links {
			urlString := htmlquery.SelectAttr(link, "href")
			if isInvalid(urlString) {
				continue
			}
			url, err := convertToAbsoluteURL(t.url, urlString)
			if err != nil {
				return err
			}

			linkText := htmlquery.InnerText(link)

			task := PageParseTask{
				url:   url,
				final: true,

				linkText: linkText,
			}
			ctx.addTask(task)
		}

	}
	return err
}

func isInvalid(url string) bool {
	result := false

	switch url {
	case "#":
		result = true
	default:
		result = false
	}

	return result
}
