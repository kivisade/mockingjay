package main

import (
	"os"
	"log"
	"net/http"
	"fmt"
	"github.com/fatih/color"
	"net/url"
	"net/http/httputil"
	"strings"
	"github.com/davecgh/go-spew/spew"
)

func echo(req *http.Request) {
	white := color.New(color.FgHiWhite).SprintfFunc()
	log.Printf("Incoming HTTP request from %s", white(req.RemoteAddr))
	log.Printf("req.RequestURI = '%s'\n", req.RequestURI)
	log.Printf("req.URL.RequestURI() = '%s'\n", req.URL.RequestURI())
	log.Printf("req.Host = '%s'\n", req.Host)
	log.Printf("req.URL.Host = '%s'\n", req.URL.Host)
	fmt.Println("--------------------------------------------------------------------------------")
	// req.Write(os.Stdout)
	if false { spew.Dump(req) }
	if dump, err := httputil.DumpRequest(req, true); err == nil {
		os.Stdout.Write(dump)
	}
	fmt.Println("\n--------------------------------------------------------------------------------")
}

func echoHandler(resp http.ResponseWriter, req *http.Request) {
	echo(req)
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	req.Write(resp)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case b == "/":
		return a
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func getProxyFunc(targetURL string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL) // target *url.URL, err error
	if err != nil {
		log.Fatal("Cannot parse URL:", targetURL)
	}
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Host = req.URL.Host
		req.RequestURI = req.URL.RequestURI()
		echo(req)
	}
	return &httputil.ReverseProxy{Director: director}
}

func main() {
	parseFlagsAndLoadConfig()

	log.Println("Starting http server on:", config.BindAddress)

	if config.ForwardTo != "" {
		log.Println("Forwarding all requests to:", config.ForwardTo)
		log.Fatal(http.ListenAndServe(config.BindAddress, getProxyFunc(config.ForwardTo)))
	} else {
		http.HandleFunc("/", echoHandler)
		log.Fatal(http.ListenAndServe(config.BindAddress, nil))
	}
}
