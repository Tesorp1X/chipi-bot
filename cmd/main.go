package main

import (
	"flag"
	"fmt"
)

var debug = flag.Bool("debug", false, "log debug info")

func main() {
	fmt.Println("hello world")
}
