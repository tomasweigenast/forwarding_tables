package main

import (
	"fmt"
	"os"
)

const endline = "\n"

type logger struct {
	f *os.File
}

var default_logger *logger = new_logger()

func new_logger() *logger {
	log_file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	log_file.Truncate(0)

	return &logger{f: log_file}
}

func (l *logger) close() {
	l.f.Close()
}

func (l *logger) log(level int, text string) {
	fmt.Println(text)
	if level == 5 {
		l.f.Write([]byte(text + endline))
	}
}

// logf prints the log message. if level is 5, it prints to file too
func (l *logger) logf(level int, text string, args ...any) {
	msg := fmt.Sprintf(text, args...)
	fmt.Println(msg)
	if level == 5 {
		l.f.Write([]byte(msg + endline))
	}
}

func (l *logger) file_logf(text string, args ...any) {
	l.logf(5, text, args...)
}

func (l *logger) infof(text string, args ...any) {
	l.logf(0, text, args...)
}
