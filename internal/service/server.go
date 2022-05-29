package service

import (
	"fmt"
	"net/http"
)

const ADDRESS = "127.0.0.1:9898"

// StartHTTPServer start new http server to grab given url when exchanging authorization code.
func StartHTTPServer(ch chan<- string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		ch <- code
		fmt.Fprint(w, "Now you may close this tab, cron-upload automatically grab the authorization code")
	})

	http.ListenAndServe(ADDRESS, nil)
}
