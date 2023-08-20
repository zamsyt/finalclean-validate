package main

import (
	"fmt"
	"log"
	"os"
)

const version = "v0.3.2"

var cmds = map[string]func(args []string){
	"fullmerge": fullMerge,
	"split":     split,
	"join":      join,
	"list":      listLayers,
	//"version":   func([]string) { fmt.Println(version) },
}

func main() {
	fmt.Printf("TFC validate %v...\n", version)
	log.SetFlags(0)
	args := os.Args[1:]
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

const basediffpath = "base-diff.png"
const palettediffpath = "palette-diff.png"

func fullMerge(args []string) {
	if len(args) == 0 {
		fmt.Println("not enough arguments")
		fmt.Println("Usage:\n\tfullmerge <my-drawpile.ora>")
		os.Exit(1)
	}
	o := OpenORA(args[0])
	f, err := o.zip.Open("mergedimage.png")
	check(err)
	mergedimg := getImg(f)
	baselayer := o.Layer("BASE LAYER")
	mergedPaletted := posterize(mergedimg, palette)
	basePaletted := posterize(baselayer, palette)

	_, baseCount := genDiff(baselayer, basePaletted)
	fmt.Println(baseCount, "pixels corrected in BASE LAYER")
	paletteDiff, mergedCount := genDiff(mergedPaletted, mergedimg)
	fmt.Println(mergedCount, "pixels corrected in the merged image")

	diff, count := genDiff(basePaletted, mergedPaletted)
	fmt.Println(count, "pixels are different between BASE LAYER and the merged image")
	fmt.Println("Saving (corrected) pixels that are different from BASE LAYER in", basediffpath)
	savePng(diff, basediffpath)

	fmt.Println("Saving non-palette pixels in", palettediffpath)
	savePng(paletteDiff, palettediffpath)
}

func printUsage() {
	fmt.Println("Available commands:")
	for k := range cmds {
		fmt.Println("\t" + k)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
