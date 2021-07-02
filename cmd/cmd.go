package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonejack/html-to-webarchive/model"
)

type HTML2WebArchive struct {
	MediaDir string
	Verbose  bool
}

func (h *HTML2WebArchive) Run(htmls []string) (err error) {
	if len(htmls) == 0 {
		return errors.New("no HTML files given")
	}

	for _, html := range htmls {
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

	return
}
