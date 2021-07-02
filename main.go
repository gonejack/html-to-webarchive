package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/gonejack/html-to-webarchive/cmd"
)

var (
	verbose bool

	prog = &cobra.Command{
		Use:   "html-to-webarchive *.html",
		Short: "Command line tool for converting html to safari .webarchive files",
		Run: func(c *cobra.Command, args []string) {
			err := run(c, args)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
)

func init() {
	log.SetOutput(os.Stdout)
	prog.Flags().SortFlags = false
	flags := prog.PersistentFlags()
	{
		flags.SortFlags = false
		flags.BoolVarP(&verbose, "verbose", "v", false, "verbose")
	}
}
func run(c *cobra.Command, args []string) error {
	exec := cmd.HTML2WebArchive{
		MediaDir: "media",
		Verbose:  verbose,
	}

	if len(args) == 0 {
		args, _ = filepath.Glob("*.html")
	}

	return exec.Run(args)
}
func main() {
	_ = prog.Execute()
}
