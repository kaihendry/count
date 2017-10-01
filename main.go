package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	jsonl "github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
)

type viewCount struct {
	sync.Mutex
	count int64
	time  time.Time
}

var v viewCount

func (n *viewCount) inc() {
	n.Lock()
	n.count++
	log.WithFields(log.Fields{
		"count": v.count,
	}).Info("count")
	n.time = time.Now()
	n.Unlock()
}

func (n *viewCount) json() []byte {
	n.time = time.Now()
	bytes, err := json.Marshal(struct {
		Count int64 `json:"count"`
		Epoch int64 `json:"time"`
	}{n.count, n.time.Unix()})
	if err != nil {
		panic(err)
	}
	return bytes
}

func init() {
	if os.Getenv("UP_STAGE") == "" {
		log.SetHandler(text.Default)
	} else {
		log.SetHandler(jsonl.Default)
	}
}

func inc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	v.inc()
	w.Write(v.json())
}

func main() {
	flag.Parse()

	// fmt.Println("ViewCount", v)
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

	t, err := template.New("foo").Parse(`<!DOCTYPE html>
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
</html>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.inc()

	envmap := make(map[string]string)
	for _, e := range os.Environ() {
		ep := strings.SplitN(e, "=", 2)
		if strings.HasPrefix(ep[0], "AWS") {
			continue
		}
		envmap[ep[0]] = ep[1]
	}
	envmap["REMOTE_ADDR"] = r.RemoteAddr
	referer := r.Referer()
	if referer != "" {
		envmap["REFERER"] = referer
	}

	err = t.Execute(w, struct {
		Count  int64
		Env    map[string]string
		Header http.Header
	}{v.count, envmap, r.Header})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
