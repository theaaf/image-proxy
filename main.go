package main

import (
	"log"
	"net/http"

	"github.aaf.cloud/platform/image-proxy/proxy"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	pr, err := proxy.NewRequestFromURL(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	proxy.Proxy(w, pr)
}

func main() {
	s := &http.Server{
		Addr:    ":8080",
		Handler: http.HandlerFunc(httpHandler),
	}
	log.Print("listenting at http://127.0.0.1:8080")
	log.Fatal(s.ListenAndServe())
}
