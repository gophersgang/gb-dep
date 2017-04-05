package cleanvcs

import (
	"flag"
	"fmt"
	"log"

	"os"

	"path/filepath"

	"github.com/gophersgang/gb-dep/pkg/gitdep"
	"github.com/gophersgang/gb-dep/pkg/packagefile"
	"github.com/gophersgang/gb-dep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("cleanvcs", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	fmt.Println("Removing VSC folders from vendor...")
	removevcs(args)
}

func (r *cmd) Usage() string {
	return "cleanvcs --verbose=true"
}

func removevcs(args []string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := packagefile.FindPackagefile(currDir)
	if err != nil {
		return err
	}
	root := filepath.Dir(file)
	pkgs, err := packagefile.Parse(file)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		d := gitdep.Dep{Pkg: pkg, RootFolder: root}
		d.Removevcs()
	}
	return nil
}
