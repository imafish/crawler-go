package main

import iconv "github.com/djimenez/iconv-go"

type groupContext struct {
	i          *counter
	dir        string
	name       string
	firstParse bool
	converter  *iconv.Converter
}
