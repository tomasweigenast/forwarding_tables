package main

import (
	"fmt"

	"github.com/arriqaaq/aol"
)

const endline = "\n"

type logger struct {
	f *aol.Log
}

func new_logger() *logger {
	log_file, err := aol.Open("log.txt", aol.DefaultOptions)
	if err != nil {
		panic(err)
	}

	return &logger{f: log_file}
}

func (l *logger) close() {
	l.f.Close()
}

func (l *logger) log(text string) {
	l.f.Write([]byte(text))
	fmt.Println(text)
}

func (l *logger) logf(text string, args ...any) {
	msg := fmt.Sprintf(text, args...)
	l.f.Write([]byte(msg))
	fmt.Println(msg)
}
