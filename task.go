package main

// Task represents a crawler rule
type Task interface {
	Execute(ctx *ExecutionContext) error
	URL() string
}
