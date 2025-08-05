package main

import (
	"cat/project"
	"fmt"
	"os"
)



func main() {
	var port string
	if len(os.Args) < 2 {
		port = ":9999"
	} else if len(os.Args) == 2 {
		port = ":" + os.Args[1]
	} else{
		fmt.Println("Usage: go run main.go [port]")
		return
	}
	project.Start(port)
}