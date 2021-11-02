package main

import (
	"log"

	"github.com/gonejack/html-to-webarchive/cmd"
)

func main() {
	c := cmd.HTMLToWarc{MediaDir: "media"}
	if e := c.Run(); e != nil {
		log.Fatal(e)
	}
}
