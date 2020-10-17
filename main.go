package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	dir, file := filepath.Split(r.URL.Path)
	// special path
	switch file {
	case "index.html":
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
	case "":
		// file is empty, check dir
		if dir == "/" {
			serveIndex(w, r)
			return
		}

		// TODO check this
		http.ServeFile(w, r, filepath.Join(".", filepath.Join(".", r.URL.Path[1:])))
	default:
		// other file names
		// check file exist
		filePath := filepath.Join(".", r.URL.Path[1:])
		fi, err := os.Stat(filePath)
		if err != nil && !os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// if filepath not exist
		if os.IsNotExist(err) && *flagRouter == "on" {
			serveIndex(w, r)
			return
		}

		// if filepath exist
		// it is dir
		if fi != nil && fi.IsDir() {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusSeeOther)
			return
		}
		// it is a file or default
		http.ServeFile(w, r, filepath.Join(".", filePath))
		return
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	indexExist, err := checkDirExists("static/index.html")
	if err != nil {
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
