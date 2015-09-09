package main

import (
	"log"
	"os"

	"github.com/yannk/scutil-to-json/scutil"
)

func main() {
	err := scutil.JSONEncode(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
