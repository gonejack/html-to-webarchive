package htm2war

import (
	"path/filepath"

	"github.com/alecthomas/kong"
)

type Options struct {
	Verbose  bool   `short:"v" help:"Verbose printing."`
	About    bool   `help:"Show about."`
	MediaDir string `hidden:"" default:"media"`

	HTML []string `name:".html" arg:"" help:"list of .html files" optional:""`
}

func MustParseOption() (opt Options) {
	kong.Parse(&opt,
		kong.Name("html-to-webarchive"),
		kong.Description("This command line converts .html to Safari's .webarchive file"),
		kong.UsageOnError(),
	)
	if len(opt.HTML) == 0 || opt.HTML[0] == "*.html" {
		opt.HTML, _ = filepath.Glob("*.html")
	}
	return
}
