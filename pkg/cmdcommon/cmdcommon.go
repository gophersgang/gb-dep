package cmdcommon

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/semaphore"
)

var (
	cfg = config.Config
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

// AllVendoredVCSFolders gives you all VCS folders
func AllVendoredVCSFolders() ([]string, error) {
	var folders []string
	searchDir, _ := cfg.AbsVendorFolder()

	err := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if isVCSFolder(path) {
			folders = append(folders, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	for _, folder := range folders {
		fmt.Println(folder)
	}
	return folders, nil
}

func isVCSFolder(path string) bool {
	if filepath.Ext(path) == ".git" || filepath.Ext(path) == ".hg" || filepath.Ext(path) == ".bzr" {
		return true
	}
	return false
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
