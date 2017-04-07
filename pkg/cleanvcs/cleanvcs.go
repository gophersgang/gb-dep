package cleanvcs

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/subcommands"
	"github.com/vbauerster/mpb"
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
	r.fs = flag.NewFlagSet("cleanvcs", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	cfg.LoggerBackend.SetPrefix("cleanvcs ")
	r.fs.Parse(args)
	fmt.Println("Removing VSC folders from vendor...")
	removevcs(args)
}

func (r *cmd) Usage() string {
	return "cleanvcs --verbose=true"
}

func removevcs(args []string) error {
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
			a.Run()
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

func bars() {
	decor := func(s *mpb.Statistics, myWidth chan<- int, maxWidth <-chan int) string {
		str := fmt.Sprintf("%3d/%3d", s.Current, s.Total)
		// send width to Progress' goroutine
		myWidth <- utf8.RuneCountInString(str)
		// receive max width
		max := <-maxWidth
		return fmt.Sprintf(fmt.Sprintf("%%%ds", max+1), str)
	}

	totalItem := 100
	var wg sync.WaitGroup
	p := mpb.New()
	wg.Add(3) // add wg delta
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Bar#%d:", i)
		bar := p.AddBar(int64(totalItem)).
			PrependName(name, len(name), 0).
			PrependFunc(decor)
		go func() {
			defer wg.Done()
			for i := 0; i < totalItem; i++ {
				bar.Incr(1)
				time.Sleep(time.Duration(rand.Intn(totalItem)) * time.Millisecond)
			}
		}()
	}
	wg.Wait() // Wait for goroutines to finish
	p.Stop()  // Stop mpb's rendering goroutine
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
