package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/skratchdot/open-golang/open"
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

func hostname() string {
	hostname, _ := os.Hostname()
	// If hostname does not have dots (i.e. not fully qualified), then return zeroconf address for LAN browsing
	if strings.Split(hostname, ".")[0] == hostname {
		return hostname + ".local"
	}
	return hostname
}

func inc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	v.inc()
	log.Println("Counter is now", v.count)
	w.Write(v.json())
}

var port = flag.Int("port", 0, "listen port")
var openbrowser = flag.Bool("openbrowser", true, "Open in browser")

func main() {
	ch := make(chan os.Signal)
	flag.Parse()

	// fmt.Println("ViewCount", v)
	http.HandleFunc("/favicon.ico", http.NotFound)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", countpage)

	// This should trigger a restart with count.service
	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ch <- syscall.SIGTERM
	})

	http.HandleFunc("/inc/", inc)

	// http://stackoverflow.com/a/33985208/4534
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panic(err)
	}

	if a, ok := ln.Addr().(*net.TCPAddr); ok {
		host := fmt.Sprintf("http://%s:%d", hostname(), a.Port)
		log.Println("Serving from", host)
		if *openbrowser {
			open.Start(host)
		}
	}

	go func() {
		if err := http.Serve(ln, nil); err != nil {
			log.Panic(err)
		}
	}()

	signal.Notify(ch, syscall.SIGTERM)
	log.Printf("Received signal '%v'. Exiting.", <-ch)

}

func countpage(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("foo").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
<meta name=viewport content="width=device-width, initial-scale=2">
<script src="static/main.js"></script>
</head>
<body>
<button onClick="f(this)">{{ .Count }}</button>

<dl>
{{range $key, $value := .Env}}
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
		ep := strings.Split(e, "=")
		envmap[ep[0]] = ep[1]
	}

	err = t.Execute(w, struct {
		Count int64
		Env   map[string]string
	}{v.count, envmap})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())

}
