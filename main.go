package main

import (
	"log"

	"github.com/gonejack/html-to-webarchive/htm2war"
)

func main() {
	cmd := htm2war.HTMLToWarc{
		Options: htm2war.MustParseOption(),
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
