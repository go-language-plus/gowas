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
	"github.com/julienschmidt/httprouter"
)

var (
	flagListen = flag.String("listen", ":5000", "HTTP listen address")
	flagDir    = flag.String("dir", ".", "directory to serve")
	flagTags   = flag.String("tags", "", "go build tags")
	flagRouter = flag.String("router", "", "router modules; if use router in wasm framework (like vecty-router), set it to 'on' ")
)

var staticOutputDir = "static/"

func main() {
	flag.Parse()
	serveHTTP()
}

func serveHTTP() {
	log.Printf("start listening http request on %q", *flagListen)

	router := httprouter.New()
	router.GET("/", handlerIndex)
	router.NotFound = http.HandlerFunc(handleNotFound)

	log.Fatal(http.ListenAndServe(*flagListen, router))
}

func handlerIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	serveStatic(w, r)
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.Redirect(w, r, r.URL.Path[:len(r.URL.Path)-1], http.StatusSeeOther)
		return
	}

	serveStatic(w, r)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {

	switch filepath.Base(r.URL.Path) {
	case "index.html", ".":
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
		fmt.Println(filepath.Join(".", r.URL.Path))
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
