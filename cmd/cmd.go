package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
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
		log.Printf("processing %s", html)

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
		log.Printf("%s exist, skipped", target)
		return
	}

	htm, err := ioutil.ReadFile(html)
	if err != nil {
		err = fmt.Errorf("open html %s fialed: %s", htm, err)
		return
	}
	wa := model.NewWebArchive(htm)

	// save resources
	err = wa.Download(html, h.MediaDir, h.Verbose)
	if err != nil {
		return
	}

	// write webarchive
	err = wa.Write(target)
	if err != nil {
		return
	}

	return
}
