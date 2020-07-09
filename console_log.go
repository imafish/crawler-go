package main

import (
	"fmt"
	"io"
)

// ConsoleLog writes logs to console
type ConsoleLog struct {
	w io.Writer
}

// Tracef trace level log
func (c ConsoleLog) Tracef(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)
}

// Debugf trace level log
func (c ConsoleLog) Debugf(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)

}

// Infof trace level log
func (c ConsoleLog) Infof(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)

}

// Warningf trace level log
func (c ConsoleLog) Warningf(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)

}

// Errorf trace level log
func (c ConsoleLog) Errorf(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)

}

// Fatalf trace level log
func (c ConsoleLog) Fatalf(format string, a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintf(c.w, format+"\n", a...)
	}
	return fmt.Printf(format+"\n", a...)

}

// Trace trace level log
func (c ConsoleLog) Trace(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// Debug trace level log
func (c ConsoleLog) Debug(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// Info trace level log
func (c ConsoleLog) Info(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// Warning trace level log
func (c ConsoleLog) Warning(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// Error trace level log
func (c ConsoleLog) Error(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// Fatal trace level log
func (c ConsoleLog) Fatal(a ...interface{}) (n int, err error) {
	if c.w != nil {
		fmt.Fprintln(c.w, a...)
	}
	return fmt.Println(a...)
}

// SetPrefix trace level log
func (c ConsoleLog) SetPrefix(prefix string) error {
	return nil
}
