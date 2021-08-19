package model

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gonejack/get"
	"howett.net/plist"
)

type Resources struct {
	WebResourceMIMEType         string `plist:"WebResourceMIMEType"`
	WebResourceTextEncodingName string `plist:"WebResourceTextEncodingName"`
	WebResourceURL              string `plist:"WebResourceURL"`
	WebResourceFrameName        string `plist:"WebResourceFrameName"`
	WebResourceData             []byte `plist:"WebResourceData"`
	//WebResourceResponse         []byte `plist:"WebResourceResponse"`
}

type WebArchive struct {
	WebMainResources *Resources   `plist:"WebMainResource"`
	WebSubResources  []*Resources `plist:"WebSubresources"`

	html string `plist:"-"`
}

func (w *WebArchive) SaveRefs(dir string, verbose bool) (err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(w.WebMainResources.WebResourceData))
	if err != nil {
		err = fmt.Errorf("parse html %s fialed: %s", w.WebMainResources.WebResourceData, err)
		return
	}

	savedFiles := make(map[string]string)
	tasks := get.NewDownloadTasks()

	doc.Find("img,video,link").Each(func(i int, e *goquery.Selection) {
		var attr string
		switch e.Get(0).Data {
		case "link":
			attr = "href"
		case "img":
			attr = "src"
			e.RemoveAttr("loading")
			e.RemoveAttr("srcset")
		case "video":
			attr = "src"
		default:
			attr = "src"
		}

		ref, _ := e.Attr(attr)
		switch {
		case ref == "":
			return
		case strings.HasPrefix(ref, "http"):
			u, err := url.Parse(ref)
			if err != nil {
				log.Printf("cannot parse %s", ref)
				return
			}
			localFile := filepath.Join(dir, fmt.Sprintf("%s%s", md5str(ref), filepath.Ext(u.Path)))
			tasks.Add(ref, localFile)
			savedFiles[ref] = localFile
		default:
			fd, err := w.openLocalFile(w.html, ref)
			if err == nil {
				_ = fd.Close()
				savedFiles[ref] = fd.Name()
			}
		}
	})

	if len(tasks.List) > 0 {
		err = os.MkdirAll(dir, 0766)
		if err != nil {
			return
		}
		get.Batch(tasks, 3, time.Minute*2).ForEach(func(t *get.DownloadTask) {
			log.Printf("download %s fail: %s", t.Link, t.Err)
		})
	}

	for ref, file := range savedFiles {
		err = w.AttachResource(ref, file)
		if err != nil {
			log.Printf("cannot attach %s(%s): %s", ref, file, err)
		}
	}

	return
}
func (w *WebArchive) AttachResource(ref string, file string) (err error) {
	fmime, err := mimetype.DetectFile(file)
	if err != nil {
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	resource := &Resources{
		WebResourceMIMEType: fmime.String(),
		WebResourceURL:      ref,
		WebResourceData:     data,
	}
	w.WebSubResources = append(w.WebSubResources, resource)

	return
}
func (w *WebArchive) Write(target string) (err error) {
	fd, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	return plist.NewEncoder(fd).Encode(w)
}
func (w *WebArchive) openLocalFile(htmlFile string, ref string) (fd *os.File, err error) {
	fd, err = os.Open(ref)
	if err == nil {
		return
	}

	// compatible with evernote's exported htmls
	{
		basename := strings.TrimSuffix(htmlFile, filepath.Ext(htmlFile))
		filename := filepath.Base(ref)
		fd, err = os.Open(filepath.Join(basename+"_files", filename))
		if err == nil {
			return
		}
		fd, err = os.Open(filepath.Join(basename+".resources", filename))
		if err == nil {
			return
		}
		if strings.HasSuffix(ref, ".") {
			return w.openLocalFile(htmlFile, strings.TrimSuffix(ref, "."))
		}
	}

	return
}

func NewWebArchive(html string) (warc *WebArchive, err error) {
	htm, err := ioutil.ReadFile(html)
	if err != nil {
		err = fmt.Errorf("open html %s fialed: %s", htm, err)
		return
	}

	warc = &WebArchive{
		WebMainResources: &Resources{
			WebResourceMIMEType:         "text/html",
			WebResourceTextEncodingName: "UTF-8",
			WebResourceURL:              "",
			WebResourceFrameName:        "",
			WebResourceData:             htm,
		},
		html: html,
	}

	return
}

func md5str(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
