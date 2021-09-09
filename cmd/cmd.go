package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/kong"

	"github.com/gonejack/html-to-webarchive/model"
)

type options struct {
	Verbose bool `short:"v" help:"Verbose printing."`
	About   bool `help:"Show about."`

	HTML []string `arg:"" optional:""`
}

type HTML2WebArchive struct {
	options
	MediaDir string
}

func (h *HTML2WebArchive) Run() (err error) {
	kong.Parse(&h.options,
		kong.Name("html-to-webarchive"),
		kong.Description("Command line tool for converting html to Safari .webarchive files."),
		kong.UsageOnError(),
	)

	if h.About {
		fmt.Println("Visit https://github.com/gonejack/html-to-webarchive")
		return
	}

	if len(h.HTML) == 0 {
		h.HTML, _ = filepath.Glob("*.html")
	}
	if len(h.HTML) == 0 {
		return errors.New("not .html file found")
	}

	err = os.MkdirAll(h.MediaDir, 0777)
	if err != nil {
		return fmt.Errorf("cannot make dir %s: %s", h.MediaDir, err)
	}

	for _, html := range h.HTML {
		log.Printf("process %s", html)

		err = h.process(html)
		if err != nil {
			return fmt.Errorf("process %s failed: %s", html, err)
		}
	}

	return
}
func (h *HTML2WebArchive) process(html string) (err error) {
	target := strings.TrimSuffix(html, filepath.Ext(html)) + ".webarchive"
	if s, e := os.Stat(target); e == nil && s.Size() > 0 {
		log.Printf("skipped %s", target)
		return
	}

	// new webarchive
	warc, err := model.NewWebArchive(html)
	if err != nil {
		return
	}

	// save resources
	err = warc.SaveRefs(h.MediaDir, h.Verbose)
	if err != nil {
		return
	}

	// write webarchive
	err = warc.Write(target)
	if err != nil {
		return
	}

	if h.Verbose {
		log.Printf("saved %s", target)
	}

	return
}
