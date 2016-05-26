package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
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

func (n *viewCount) getCount() int64 {
	return n.count
}

func (n *viewCount) json() []byte {
	n.count++
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
	w.Write(v.json())
}

func main() {

	fmt.Println("ViewCount", v)
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", lk)
	http.HandleFunc("/inc/", inc)

	// http://stackoverflow.com/a/33985208/4534
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", 0))
	if err != nil {
		log.Panic(err)
	}

	if a, ok := ln.Addr().(*net.TCPAddr); ok {
		host := fmt.Sprintf("http://%s:%d", hostname(), a.Port)
		log.Println("Serving from", host)
		open.Start(host)
	}
	if err := http.Serve(ln, nil); err != nil {
		log.Panic(err)
	}

}

func lk(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("foo").Parse(`<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
</head>
<body>
<h1>Count {{ .Count }}</h1>
</body>
</html>`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	v.inc()
	t.Execute(w, struct {
		Count int64
	}{v.count})

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())

}
