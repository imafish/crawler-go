package main

// Task represents a crawler rule
type Task interface {
	Execute(ctx *Context) error
	URL() string
}
