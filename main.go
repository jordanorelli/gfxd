package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

var options struct {
	Host string
}

func init() {
	flag.StringVar(&options.Host, "host", "0.0.0.0:8000", "http hostname:port to listen on")
}

func main() {
	flag.Parse()
	s := server{
		out:    os.Stdout,
		errors: os.Stderr,
	}
	if err := http.ListenAndServe(options.Host, &s); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
