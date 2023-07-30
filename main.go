package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"os"
)

var cmds = map[string]func(args []string){
	"palette": checkPalette,
}

func checkPalette(args []string) {
	if len(args) < 2 {
		log.Fatal("too few arguments")
	}
	o := OpenORA(args[1])
	f, err := o.zip.Open("mergedimage.png")
	check(err)
	diff, n := genDiff(getImg(f))
	fmt.Printf("%v pixels differ from the palette\n", n)
	savePng(diff, "output.png")
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
