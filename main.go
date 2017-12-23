package main

import (
	"flag"
	"log"
	"net/http"
)

type app struct {
	static http.Handler
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// strip trailing slash.
	if url := r.URL.Path; len(url) > 1 && url[len(url)-1] == '/' {
		http.Redirect(w, r, url[:len(url)-1], 301)
		return
	}

	w.Write([]byte("hello, world"))
}

func main() {
	assets := flag.String("assets", "./assets", "Location of asset files -- including templates and static files.")
	flag.Parse()

	server := http.Server{
		Addr:    ":9393",
		Handler: &app{},
	}

	log.Printf("assets dir: %q", *assets)

	log.Fatal(server.ListenAndServe())
}
