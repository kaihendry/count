package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/skratchdot/open-golang/open"
)

var viewCount struct {
	sync.Mutex
	Count int64
}

func hostname() string {
	hostname, _ := os.Hostname()
	// If hostname does not have dots (i.e. not fully qualified), then return zeroconf address for LAN browsing
	if strings.Split(hostname, ".")[0] == hostname {
		return hostname + ".local"
	}
	return hostname
}

func main() {

	http.HandleFunc("/", lk)

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

	fmt.Println(viewCount.Count)
	viewCount.Lock()
	viewCount.Count++
	viewCount.Unlock()
	t.Execute(w, viewCount)

	log.Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())

}
