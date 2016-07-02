package main

import (
	"fmt"
	"io"
	"net/http"
)

type server struct {
	out    io.Writer
	errors io.Writer
}

func (s *server) logReceived(r *http.Request) {
	fmt.Fprintf(s.out, "> %s\n", r.URL.String())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.logReceived(r)
}
