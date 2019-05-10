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

type viewCount int32

var v viewCount

func (n *viewCount) inc() (currentcount int32) {
	return atomic.AddInt32((*int32)(n), 1)
}

func inc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%d", v.inc()) // actual valid
}

func routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", http.NotFound)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/metrics", prometheus)
	mux.HandleFunc("/", countpage)
	mux.HandleFunc("/inc/", inc)
	return mux
}

func main() { log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), routes())) }

func countpage(w http.ResponseWriter, r *http.Request) {

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
	}{v.inc(), envmap, r.Header})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func prometheus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "# HELP count_total shows the in-memory count, which will get reset in the event of the lambda going cold or scaling.\n# TYPE count_total counter\ncount_total %d\n", v)
}
