package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-language-plus/gowas/config"
	"github.com/go-language-plus/pkg/stringp"
)

var (
	flagListen = flag.String("listen", ":5000", "HTTP listen address")
	flagDir    = flag.String("dir", ".", "directory to serve")
	flagRouter = flag.String("router", "", "router modules; if use router in wasm framework (like vecty-router), set it to 'on' ")
	flagTags   = flag.String("tags", "", "go build tags")
)

var staticOutputDir = "static/"

func main() {
	flag.Parse()
	serveHTTP()
}

func serveHTTP() {
	log.Printf("start listening http request on %q", *flagListen)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(*flagListen, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[1:]
	filePath := filepath.Join(".", filepath.Base(url))

	if !strings.HasSuffix(r.URL.Path, "/") {
		fi, err := os.Stat(filePath)
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if fi != nil && fi.IsDir() {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusSeeOther)
			return
		}
	}

	switch filepath.Base(r.URL.Path) {
	case "index.html":
		fmt.Println(r.URL.Path + "inside")
		serveIndex(w, r)
		return
	case "main.wasm":
		if stdoutStderr, err := buildWasm(); err != nil {
			http.Error(w, string(stdoutStderr), http.StatusInternalServerError)
			return
		}

		f, err := os.Open(filepath.Join(staticOutputDir, "main.wasm"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		http.ServeContent(w, r, "main.wasm", time.Now(), f)
		return
	case "wasm_exec.js":
		f := filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js")
		http.ServeFile(w, r, f)
		return
	}

	if *flagRouter == "on" {
		serveIndex(w, r)
	} else {
		http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexExist, err := checkDirExists("static/index.html")
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if indexExist {
		indexHTMl := filepath.Join(".", "static", "index.html")
		log.Printf("find index.html at '%s'", indexHTMl)
		http.ServeFile(w, r, indexHTMl)
	} else {
		// serve default setting index.html
		log.Printf("index.html does not exist in the current workspace, load default index HTML file")
		byteHtml := stringp.ByteString(config.DefaultIndexHtml).SliceByte() // string to []byte
		http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader(byteHtml))
	}
}
