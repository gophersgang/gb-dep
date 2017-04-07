package install

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("install", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	fmt.Println("Running install....")
	install(args)
}

func (r *cmd) Usage() string {
	return "install --verbose=true"
}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}

func install(args []string) error {
	pkgs, err := packagefile.GimmePackagefile()
	checkErr("Gimmi packages", err)
	pwd, err := os.Getwd()
	root, err := packagefile.RootDir(pwd)
	checkErr("-", err)
	deps := []*dep.Dep{}

	for _, pkg := range pkgs.Packages {
		a := pkg
		d := &dep.Dep{Pkg: &a, RootFolder: root}
		deps = append(deps, d)
	}

	for _, d := range deps {
		d.Run()
	}

	packagefile.GenerateLockFile(deps)
	return nil
}
