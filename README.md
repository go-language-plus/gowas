# Gowas

![Go](https://github.com/go-language-plus/gowas/workflows/Go/badge.svg?branch=main)

Gowas（go-webassembly-serve）Is a webassembly development tool that provides http service support.

This repo is inspired by: [hajimehoshi/wasmserve](https://github.com/hajimehoshi/wasmserve), but some content was changed. 
The problem that this tool wants to solve is whether different wasm encoding methods can be supported at the same time, including frameworks. 
All the go webassembly frameworks seems to be in the experimental stage, so Vecty is currently prioritized because it looks most in line with the go concept.

What does this tool do:
- Auto run go build main.wasm
- Auto load index.html (if not exist in current workspaces) and  wasm_exec.js
- All in one  command line

less/scss builder is under consideration and `this repo is still experimenting` about how to write wasm (front end) with go comfortably.

## Install
```bash
go get -u github.com/go-language-plus/gowas
```
Make sure your `$GOPATH/bin` has been set in the environment variable so that we can use the `gowas` command directly.

## Usage
Enter the root directory of your webassembly project (where the `main.go` file located).

Use：
```
gowas
```
Or:
```
gowas -router on -listen :8080
```

Other Tags:
```
gowas
    -listen 
        HTTP listen address. defult: :5000
    -dir 
        directory to serve defult: .
    -tags 
        go build tags
    -router 
        router modules; if use router in wasm framework (like vecty-router), set it to 'on'
```
