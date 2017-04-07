package buildbins

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

var (
	cfg = config.Config
)

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("buildbins", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	cfg.LoggerBackend.SetPrefix("buildbins ")
	if r.verbose {
		cfg.SetDebugMode()
	}
	cfg.Logger.Println("Build all the vendor binaries / libraries")
	build(args)
}

func (r *cmd) Usage() string {
	return "buildbins --verbose=true"
}

func build(args []string) error {
	deps()
	return nil
}

func deps() {
	s := newSemaphore(5)
	var wg sync.WaitGroup
	deps := all()
	for _, d := range deps {
		s.Acquire(1)
		wg.Add(1)
		a := d
		go func() {
			defer wg.Done()
			defer s.Release(1)
			a.BuildBins()
		}()
	}
	wg.Wait()
	packagefile.GenerateLockFile(deps)
}

func all() []*dep.Dep {
	pkgs, err := packagefile.GimmePackagefile()
	checkErr("GimmePackagefile", err)
	pwd, err := os.Getwd()
	root, err := packagefile.RootDir(pwd)
	checkErr("RootDir", err)
	deps := []*dep.Dep{}

	for _, pkg := range pkgs.Packages {
		a := pkg
		d := &dep.Dep{Pkg: &a, RootFolder: root}
		deps = append(deps, d)
	}
	return deps
}

// Semaphore controls access to a finite number of resources.
type Semaphore chan struct{}

// New creates a Semaphore that controls access to `n` resources.
func newSemaphore(n int) Semaphore {
	return Semaphore(make(chan struct{}, n))
}

// Acquire `n` resources.
func (s Semaphore) Acquire(n int) {
	for i := 0; i < n; i++ {
		s <- struct{}{}
	}
}

// Release `n` resources.
func (s Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}
