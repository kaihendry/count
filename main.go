package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/apex/log"
	jsonl "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
	"github.com/tj/go/http/response"
)

type viewCount struct {
	sync.RWMutex
	count int
}

var v viewCount

func (n *viewCount) inc() (currentcount int) {
	n.Lock()
	// atomic.AddInt32(n.count, 1)
	n.count++
	currentcount = n.count
	n.Unlock()
	return currentcount
}

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonl.Default)
	}
}

func inc(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, v.inc())
}

func main() {
	flag.Parse()

	http.HandleFunc("/favicon.ico", http.NotFound)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", countpage)

	http.HandleFunc("/inc/", inc)
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatalf("error listening: %s", err)
	}

}

func countpage(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang=en>
<head>
<meta charset="utf-8" />
<meta name=viewport content="width=device-width, initial-scale=1">
<script src="static/main.js"></script>
<title>Count: {{ .Count }}</title>
<style>
body { background-color: pink; font-family: Georgia; }
</style>
</head>
<body>
<button onClick="f(this)">{{ .Count }}</button>

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
</html>`))

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

	err := t.Execute(w, struct {
		Count  int
		Env    map[string]string
		Header http.Header
	}{v.inc(), envmap, r.Header})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
