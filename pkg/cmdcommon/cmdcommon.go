package cmdcommon

import (
	"fmt"
	"os"
	"sync"

	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/semaphore"
)

// CurrentDeps returns packages wrapped in dep structs
func CurrentDeps() []*dep.Dep {
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

// RunConcurrently executes in concurrent fashion a function over given dep structs
func RunConcurrently(deps []*dep.Dep, limit int, action func(d *dep.Dep)) {
	s := semaphore.New(limit)
	var wg sync.WaitGroup
	for _, d := range deps {
		s.Acquire(1)
		wg.Add(1)
		a := d
		go func() {
			defer wg.Done()
			defer s.Release(1)
			action(a)
		}()
	}
	wg.Wait()
}

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
		os.Exit(1)
	}
}
