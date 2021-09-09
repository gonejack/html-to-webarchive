package main

import (
	"log"

	"github.com/gonejack/html-to-webarchive/cmd"
)

func main() {
	c := cmd.HTML2WebArchive{
		MediaDir: "media",
	}
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
}
