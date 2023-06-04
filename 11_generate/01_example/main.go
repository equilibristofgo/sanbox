package main

import "fmt"

//go:generate ./command.sh

func main() {
	fmt.Println("if you type 'go generate' in this directory command.sh will be run")
}
