package main

import "os"

func main() {
	if os.Args[1] == "extract" {
		extract()
	} else if os.Args[1] == "render" {
		render()
	} else {
		os.Exit(1)
	}
}
