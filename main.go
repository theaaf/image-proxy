package main

import (
	"log"
	"net/http"

	"github.com/theaaf/image-proxy/proxy"
)

func httpHandler(config *proxy.Configuration) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pr, err := proxy.NewRequestFromURL(r.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		proxy.Proxy(config, w, pr)
	})
}

func main() {
	config := &proxy.Configuration{}
	config.LoadEnvironmentVariables()

	s := &http.Server{
		Addr:    ":8080",
		Handler: httpHandler(config),
	}
	log.Print("listenting at http://127.0.0.1:8080")
	log.Fatal(s.ListenAndServe())
}
