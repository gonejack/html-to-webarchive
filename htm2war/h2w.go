package htm2war

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gonejack/html-to-webarchive/model"
)

type HTMLToWarc struct {
	Options
}

func (h *HTMLToWarc) Run() (err error) {
	if h.About {
		fmt.Println("Visit https://github.com/gonejack/html-to-webarchive")
		return
	}
	if len(h.HTML) == 0 {
		return errors.New("no .html files found")
	}
	return h.run()
}
func (h *HTMLToWarc) run() (err error) {
	for _, html := range h.HTML {
		log.Printf("process %s", html)

		err = h.process(html)
		if err != nil {
			return fmt.Errorf("process %s failed: %s", html, err)
		}
	}
	return
}
func (h *HTMLToWarc) process(html string) (err error) {
	output := strings.TrimSuffix(html, filepath.Ext(html)) + ".webarchive"
	if s, e := os.Stat(output); e == nil && s.Size() > 0 {
		log.Printf("skipped %s", output)
		return
	}

	// new webarchive
	warc, err := model.NewWebArchive(html)
	if err != nil {
		return
	}

	// save resources
	_ = os.MkdirAll(h.MediaDir, 0766)
	err = warc.SaveRefs(h.MediaDir, h.Verbose)
	if err != nil {
		return
	}

	// write webarchive
	err = warc.Write(output)
	if err != nil {
		return
	}

	if h.Verbose {
		log.Printf("saved %s", output)
	}

	return
}
