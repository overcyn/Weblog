// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"net/http"
	"strings"
	"time"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	logs          = []Log{}
)

type indexParams struct {
	Logs []Log
}

func main() {
	http.HandleFunc("/", handle)
	_ = http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
		indexTemplate.Execute(w, indexParams{Logs: logs})
	case http.MethodPost:
		// Create a new record.
		t, err := time.Parse(time.RFC3339, r.FormValue("date"))
		if err != nil {
			t = time.Now()
		}

		l := Log{
			Message: r.FormValue("message"),
			Date:    t,
			Tags:    strings.Split(r.FormValue("tags"), ","),
		}
		logs = append(logs, l)
	case http.MethodDelete:
		logs = []Log{}
	default:
		// Give an error message.
		http.Error(w, "Invalid request method.", 405)
	}
}

type Log struct {
	Message string
	Date    time.Time
	Tags    []string
}
