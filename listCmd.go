package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var listset = pflag.NewFlagSet("list", pflag.ExitOnError)
var layerNumFlag = listset.BoolP("layer-number", "n", false, "print layer numbers")

func listLayers(args []string) {
	listset.Parse(args)
	args = listset.Args()
	if len(args) == 0 {
		printListUsage()
	}
	o := OpenORA(args[0])
	for _, l := range o.Layers {
		if *layerNumFlag {
			fmt.Printf("%v ", l.Src[len("data/layer"):len(l.Src)-len(".png")])
		}
		fmt.Printf("%q\n", l.Name)
	}
}

func printListUsage() {
	fmt.Println("not enough arguments")
	fmt.Printf("Usage:\n\tlist <my-drawpile.ora>\n\n")
	listset.PrintDefaults()
	os.Exit(1)
}
