package main

import "syscall/js"

var document = js.Global().Get("document")

func main() {
	// create a title
	h1 := document.Call("createElement", "H1")
	h1.Set("innerText", "Go in frontend!")

	// set this title into <body>
	document.Get("body").Call("appendChild", h1)
}
