package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

type countHandler struct{ n int32 }

func main() { log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), routes())) }

func (h *countHandler) inc() int32 {
	return atomic.AddInt32(&h.n, 1)
}

func routes() *http.ServeMux {
	h := countHandler{}
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/favicon.ico", http.NotFound)
	mux.HandleFunc("/metrics", h.prometheus)
	mux.HandleFunc("/", h.countpage)
	mux.HandleFunc("/inc/", h.json)
	return mux
}

func (h *countHandler) json(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", h.inc()) // actually valid JSON
}

func (h *countHandler) countpage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").ParseFiles("static/index.tmpl"))
	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		// Skip potentially security sensitive AWS stuff
		if ep[0] == "AWS_SECRET_ACCESS_KEY" {
			continue
		}
		if ep[0] == "AWS_SESSION_TOKEN" {
			continue
		}

		envmap[ep[0]] = ep[1]
	}

	// https://golang.org/pkg/net/http/#Request
	envmap["METHOD"] = r.Method
	envmap["PROTO"] = r.Proto
	envmap["CONTENTLENGTH"] = fmt.Sprintf("%d", r.ContentLength)
	envmap["TRANSFERENCODING"] = strings.Join(r.TransferEncoding, ",")
	envmap["REMOTEADDR"] = r.RemoteAddr
	envmap["HOST"] = r.Host
	envmap["REQUESTURI"] = r.RequestURI

	err := t.ExecuteTemplate(w, "index.tmpl", struct {
		Count  int32
		Env    map[string]string
		Header http.Header
	}{h.inc(), envmap, r.Header})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *countHandler) prometheus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `# HELP count_total shows the in-memory count, which will get reset in the event of the lambda going cold or scaling.
# TYPE count_total counter
count_total %d`, h.n)
}
