package main

import (
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type app struct {
	static   http.Handler
	log      *log.Logger
	template *template.Template
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// strip trailing slash.
	if url := r.URL.Path; len(url) > 1 && url[len(url)-1] == '/' {
		http.Redirect(w, r, url[:len(url)-1], 301)
		return
	}

	a.log.Printf("request for %q", r.URL.Path)

	switch {
	case strings.HasPrefix(r.URL.Path, "/static/"):
		a.static.ServeHTTP(w, r)
	default:
		a.template.ExecuteTemplate(w, "index.html", struct{ Title string }{Title: "The Last Jedi Sucks!"})
	}

}

func main() {
	addr := flag.String("addr", ":9393", "Address to listen on.")
	assets := flag.String("assets", "./assets", "Location of asset files -- including templates and static files.")
	quiet := flag.Bool("quiet", false, "Disable log output.")
	flag.Parse()

	logWriter := func() io.Writer {
		if *quiet == true {
			return ioutil.Discard
		}
		return os.Stderr
	}()

	tmpl, err := template.ParseGlob(path.Join(*assets, "template/*"))

	if err != nil {
		log.Fatal(err)
	}

	tljs := app{
		log:      log.New(logWriter, "DEBUG ", log.LstdFlags|log.Lshortfile),
		template: tmpl,
		static:   http.FileServer(http.Dir(*assets)),
	}

	server := http.Server{
		Addr:    *addr,
		Handler: &tljs,
	}

	tljs.log.Printf("assets dir %q", *assets)
	tljs.log.Printf("lisening at %q", *addr)

	log.Fatal(server.ListenAndServe())
}
