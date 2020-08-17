package main

import (
	"fmt"

	"github.com/antchfx/htmlquery"
	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/net/html"
)

// PageParseTask holds information to parse a page
type PageParseTask struct {
	url         string
	final       bool
	newGroup    bool
	groupFormat *group

	taskContext *TaskContext
}

// URL returns the task's corresponding url
func (t PageParseTask) URL() string {
	return t.url
}

// Execute grab and analyse the page content, generates other related tasks.
func (t PageParseTask) Execute(ctx *ExecutionContext) error {
	reader, err := grabPageReader(t.url)
	if err != nil {
		GetLogger().Errorf("failed to grab page from url %s, err: %s", t.url, err.Error())
		return err
	}
	doc, err := htmlquery.Parse(reader)
	if err != nil {
		GetLogger().Errorf("failed to parse page content. err: %s", err.Error())
		return err
	}

	encoding := GetPageEncoding(doc)
	var converter *iconv.Converter
	if encoding != "" && encoding != "utf-8" {
		converter, err = iconv.NewConverter(encoding, "utf-8")
		if err != nil {
			GetLogger().Errorf("failed to create encoding converter from %s. err: %s", encoding, err.Error())
			return err
		}
	}

	// fill in task context
	t.taskContext.pageURL = t.url
	t.taskContext.pageTitle = GetPageTitle(doc, converter)
	if t.newGroup {
		t.taskContext.startPageTitle = t.taskContext.pageTitle
		t.taskContext.startPageURL = t.taskContext.pageURL
	}

	groupContext := ctx.FindOrNewGroupContext(t.taskContext, t.newGroup, t.groupFormat)
	if groupContext != nil {
		t.taskContext.counter = groupContext.counter
	}

	matchingRules := ctx.FindMatchingRules(t.taskContext)

	for _, r := range matchingRules {
		for _, a := range r.Actions {
			if a.ProcessLink != nil && t.final {
				continue
			}

			matchingNodes := make([]*html.Node, 0)
			for _, tgt := range a.Targets {
				m, err := GetMatchingNodes(doc, tgt)
				if err != nil {
					GetLogger().Errorf("Error in getting matching nodes. err: %s", err.Error())
				}
				matchingNodes = append(matchingNodes, m...)
			}

			// TODO @imafish refactor this code
			for _, n := range matchingNodes {
				GetLogger().Debugf("Processing a matching node: %v", n)

				if df := a.DownloadFile; df != nil {
					url, err := t.getURL(n, df.Target)
					if err != nil {
						GetLogger().Warningf("Failed to get URL (target) from node %v", n)
						continue
					}
					taskContext, err := t.createTaskContext(url, n, converter)
					if err != nil {
						GetLogger().Warningf("Error when trying to create task context, node: %v", n)
						continue
					}

					task := DownloadFileTask{
						url:             url,
						dirPattern:      df.DirPattern,
						filenamePattern: df.FilenamePattern,
						taskContext:     taskContext,
					}
					ctx.AddTask(task)

				} else if gt := a.GrabText; gt != nil {
					url := t.url
					taskContext, err := t.createTaskContext(url, n, converter)
					if err != nil {
						GetLogger().Warningf("Error when trying to create task context, node: %v", n)
						continue
					}

					task := GrabTextTask{
						taskContext:     taskContext,
						node:            n,
						dirPattern:      gt.DirPattern,
						filenamePattern: gt.FilenamePattern,
						target:          gt.Target,
					}
					ctx.AddTask(task)

				} else if pl := a.ProcessLink; pl != nil {
					url, err := t.getURL(n, pl.Target)
					if err != nil {
						GetLogger().Warningf("Failed to get URL (target) from node %v", n)
						continue
					}

					taskContext, err := t.createTaskContext(url, n, converter)
					if err != nil {
						GetLogger().Warning("Error when trying to create task context, node: %v", n)
						continue
					}

					task := PageParseTask{
						url:         url,
						final:       pl.Final,
						newGroup:    pl.NewGroup,
						groupFormat: pl.Group,
						taskContext: taskContext,
					}
					ctx.AddTask(task)

				} else {
					GetLogger().Warning("Action does not have anything to do.")
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

func (t PageParseTask) getURL(n *html.Node, target target) (string, error) {
	urlString, err := getOneNodeInnerText(n, target)
	if err != nil {
		GetLogger().Errorf("Error in getting inner text of sub-target, err: %s", err.Error())
		return "", err
	}
	if isInvalid(urlString) {
		GetLogger().Warningf("Got invalid URL: %s", urlString)
		return "", fmt.Errorf("Got an invalid URL: %s", urlString)
	}
	urlString, err = convertToAbsoluteURL(t.url, urlString)
	if err != nil {
		GetLogger().Errorf("Unable to convert to absolute path, err: %s", err.Error())
		return "", err
	}

	return urlString, nil
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

func (t PageParseTask) createTaskContext(url string, n *html.Node, converter *iconv.Converter) (*TaskContext, error) {
	imgAlt := htmlquery.SelectAttr(n, "alt")
	imgAlt = cc(imgAlt, converter)

	ext, err := GetExtension(url)
	if err != nil {
		GetLogger().Warningf("Error in getting extension from URL: %s, error is %s", url, err.Error())
	}

	taskContext := &TaskContext{
		startPageURL:   t.taskContext.startPageURL,
		startPageTitle: t.taskContext.startPageTitle,
		pageURL:        t.taskContext.pageURL,
		pageTitle:      t.taskContext.pageTitle,
		linkURL:        url,
		linkText:       CollectText(n, converter),
		imgAlt:         imgAlt,
		extension:      ext,
		counter:        t.taskContext.counter,
	}

	return taskContext, nil
}
