package main

import (
	"fmt"
	"os"
)

var cmds = map[string]func(args []string){
	"hello": func(args []string) { fmt.Println("hi") },
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}
	cmd, ok := cmds[os.Args[1]]
	if !ok {
		fmt.Println("Unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
	cmd(os.Args[1:])
}

func printUsage() {
	fmt.Println("Available commands:")
	for k := range cmds {
		fmt.Println("\t" + k)
	}
}
