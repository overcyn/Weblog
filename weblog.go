// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"html/template"
	"net/http"
	"time"
)

var (
	indexTemplate = template.Must(template.ParseFiles("index.html"))
)

type indexParams struct {
	Logs []Log
}

func main() {
	http.HandleFunc("/", handle)
	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	switch r.Method {
	case http.MethodGet:
		// Serve the resource.
		var logs []Log
		q := datastore.NewQuery("Log").Order("-Date").Limit(1000)
		if _, err := q.GetAll(ctx, &logs); err != nil {
			log.Errorf(ctx, "datastore.GetAll: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		indexTemplate.Execute(w, indexParams{Logs: logs})
	case http.MethodPost:
		// Create a new record.
		key := datastore.NewIncompleteKey(ctx, "Log", nil)
		l := Log{
			Message: r.FormValue("message"),
			Date:    time.Now(),
		}
		if _, err := datastore.Put(ctx, key, &l); err != nil {
			log.Errorf(ctx, "datastore.Put: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodDelete:
		// Remove the record.
		// Serve the resource.
		q := datastore.NewQuery("Log").Limit(1000).KeysOnly()
		keys, err := q.GetAll(ctx, nil)
		if err != nil {
			log.Errorf(ctx, "datastore.GetAll: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := datastore.DeleteMulti(ctx, keys); err != nil {
			log.Errorf(ctx, "datastore.DeleteMulti: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		// Give an error message.
		http.Error(w, "Invalid request method.", 405)
	}

}

type Log struct {
	Message string
	Date    time.Time
}
