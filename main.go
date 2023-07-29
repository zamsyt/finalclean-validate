package main

import (
	"fmt"
	"log"
	"os"
)

func testing(args []string) {
	if len(args) < 2 {
		log.Fatal("too few arguments")
	}
	o := OpenORA(args[1])
	img := o.Layer("BASE LAYER")
	fmt.Println(img.Bounds())
	check(o.Close())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var cmds = map[string]func(args []string){
	"testing": testing,
}

func main() {
	log.SetFlags(0)
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
