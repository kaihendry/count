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

var Version string

type countHandler struct{ n int32 }

func main() { log.Fatal(http.ListenAndServe(":"+os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT"), routes())) }

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

	// Use go:embed https://github.com/golang/go/issues/41191 once it comes out
	t, err := template.New("index").Parse(`
<!DOCTYPE html>
<html lang=en>
<head>
<meta charset="utf-8" />
<meta name=viewport content="width=device-width, initial-scale=1">
<title>Count: {{ .Count }}</title>
<style>
body { background-color: white; font-family: Georgia; }
</style>
</head>
<body>
<dl>
{{range $key, $value := .Env -}}
{{ if eq $key "COMMIT" -}}
<dt>{{ $key }}</dt><dd><a href="https://github.com/kaihendry/count/commit/{{ $value }}">{{ $value }}</a></dd>
{{else}}
<dt>{{ $key }}</dt><dd>{{ $value }}</dd>
{{- end}}
{{- end}}
</dl>
<h3>Request Header</h3>
<dl>
{{range $key, $value := .Header -}}
<dt>{{ $key }}</dt><dd>{{ $value }}</dd>
{{end}}
</dl>
<p><a href=https://github.com/kaihendry/count>Source code</a></p>
</body>
</html>`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		envmap[ep[0]] = ep[1]
	}

	// https://golang.org/pkg/net/http/#Request
	envmap["COUNTVERSION"] = Version
	envmap["METHOD"] = r.Method
	envmap["PROTO"] = r.Proto
	envmap["CONTENTLENGTH"] = fmt.Sprintf("%d", r.ContentLength)
	envmap["TRANSFERENCODING"] = strings.Join(r.TransferEncoding, ",")
	envmap["REMOTEADDR"] = r.RemoteAddr
	envmap["HOST"] = r.Host
	envmap["REQUESTURI"] = r.RequestURI

	err = t.Execute(w, struct {
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
