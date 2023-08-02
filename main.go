package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
)

const version = "v0.1.0"

var cmds = map[string]func(args []string){
	"checkpalette": checkPalette,
	"fixpalette":   fixPalette,
	"version":      func([]string) { fmt.Println(version) },
}

func main() {
	pflag.Parse()
	log.SetFlags(0)
	args := pflag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(0)
	}
	cmd, ok := cmds[args[0]]
	if !ok {
		fmt.Println("Unknown command:", args[0])
		printUsage()
		os.Exit(1)
	}
	cmd(args[1:])
}

func printUsage() {
	fmt.Println("Available commands:")
	for k := range cmds {
		fmt.Println("\t" + k)
	}
}

func getImg(r io.ReadCloser) image.Image {
	img, _, err := image.Decode(r)
	check(err)
	check(r.Close())
	return img
}

func savePng(img image.Image, path string) {
	f, err := os.Create(path)
	check(err)
	png.Encode(f, img)
	f.Close()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
