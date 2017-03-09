package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"io"
	"fmt"
	"time"
	"runtime"

	"github.com/fatih/color"
)

func handler(resp http.ResponseWriter, req *http.Request) {
	t := time.Now()
	white := color.New(color.FgHiWhite).SprintfFunc()
	fmt.Println("\n================================================================================")
	fmt.Printf("Request: %s\n", white(t.Format("2006-01-02 15:04:05"))) // Format(time.RFC3339)
	fmt.Println("--------------------------------------------------------------------------------")

	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	req.Write(io.MultiWriter(resp, os.Stdout))

	fmt.Println()
}

func main() {
	bindingAddress := flag.String("listen", ":8080", "Listen address:port (defaults to :8080)")
	flagNoColor := flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	if *flagNoColor || runtime.GOOS == "windows" {
		color.NoColor = true // disables colorized output
	}

	http.HandleFunc("/", handler)
	log.Println("Starting http server on:", *bindingAddress)
	log.Fatal(http.ListenAndServe(*bindingAddress, nil))
}
