// Copyright 2014 errnoh. All rights reserved.
// Use of this source code is governed by
// MIT License that can be found in the LICENSE file.
package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"strings"
)

var supported = make(map[string]func(*http.Request, io.Reader, *bytes.Buffer))

func main() {
	http.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.ContentLength == 0 || r.Body == nil {
		http.NotFound(w, r)
		return
	}

	functions := strings.Split(html.EscapeString(r.URL.Path), "/")
	if len(functions) <= 1 {
		http.NotFound(w, r)
		return
	}
	// Path should always start with / so we'll drop the empty string at [0]
	functions = functions[1:]

	// Check that all requested functions are supported
	for i := 0; i < len(functions); i++ {
		if _, ok := supported[functions[i]]; !ok {
			http.NotFound(w, r)
			return
		}
	}

	b := []*bytes.Buffer{new(bytes.Buffer), new(bytes.Buffer)}

	// Call the first operation with request body as source, use the piped data for rest of the operations
	supported[functions[0]](r, r.Body, b[0])
	for i := 1; i < len(functions); i++ {
		supported[functions[i]](r, b[(i-1)%2], b[i%2])
	}

	io.Copy(w, b[(len(functions)-1)%2])
}
