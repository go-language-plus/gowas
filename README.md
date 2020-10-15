# Gowas

![Go](https://github.com/go-language-plus/gowas/workflows/Go/badge.svg?branch=main)

Gowas（go-webassembly-serve）Is a Wasm development tool that provides http service support.

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

