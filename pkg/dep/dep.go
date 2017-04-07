package dep

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/gbutils"
	"github.com/gophersgang/gbdep/pkg/packagefile"
)

var (
	cfg = config.Config
)

// Dep is a package
type Dep struct {
	packagefile.Pkg
	RootFolder string // the root folder
	// the subpath for GIT / HG folder, that might not match the full package path,
	// eg. golang.org/x/crypto/ssh -> golang.org/x/crypto/.git
	VcsFolder string
}

// Run knows what to do
func (d *Dep) Run() error {
	fmt.Println("Installing " + d.Name)
	d.ensureProperEnv()
	d.ensureBasefolders()
	d.ensureInstalled()
	return nil
}

func (d *Dep) ensureInstalled() error {
	d.detectVcsFolder()
	if !d.cacheExists() {
		d.slowInstall()
		return d.Copyvcs()
	}

	return d.installFromCache()
}

// this methods does the traditional `go get` installation
// happens only if you dont have any cached VCS folder
func (d *Dep) slowInstall() error {
	plainRunCmd(d.vendorFolder(), fmt.Sprintf("go get -u %s", d.Name))
	plainRunCmd(d.pkgVendorFolder(), "go get -u ./...")
	plainRunCmd(d.pkgVendorFolder(), "go install ./...")
	return nil
}

// this is executed on consecutive runs and should be much faster
func (d *Dep) installFromCache() error {
	fmt.Println("Installing from cache...")
	return nil
}

func (d *Dep) ensureProperEnv() {
	os.Setenv("GOBIN", filepath.Join(d.vendorFolder(), "bin")) // this is important to get binaries
	os.Setenv("GOPATH", filepath.Join(d.vendorFolder()))       // so it works reliably
}

func (d *Dep) ensureBasefolders() error {
	os.MkdirAll(d.cacheFolder(), 0777)
	os.MkdirAll(d.vendorFolder(), 0777)
	os.MkdirAll(d.pkgCachedFolder(), 0777)
	return nil
}

// Copyvcs creates .GIT folder in vendor, for installation and such
// run before commiting
func (d *Dep) Copyvcs() error {
	path, err := d.detectVcsFolder()
	if err != nil {
		return err
	}
	fromPath := filepath.Join(d.vendorFolder(), "src", path)
	toPath := filepath.Join(d.cacheFolder(), path)
	// TODO: think about proper update for cached folders
	if gbutils.PathExist(toPath) {
		return nil
	}
	copyCmd := fmt.Sprintf("cp -r %s %s", fromPath, toPath)
	return plainRunCmd(d.vendorFolder(), copyCmd)
}

func (d *Dep) cacheExists() bool {
	return gbutils.PathExist(filepath.Join(d.cacheFolder(), d.VcsFolder))
}

func (d *Dep) pkgCachedFolder() string {
	return filepath.Join(d.cacheFolder(), d.Name)
}

//
func (d *Dep) pkgVendorFolder() string {
	return filepath.Join(d.vendorFolder(), "src", d.Name)
}

func (d *Dep) cacheFolder() string {
	return d.RootFolder + "/" + "vendor.cache"
}

func (d *Dep) vendorFolder() string {
	return d.RootFolder + "/" + "vendor"
}

func (d *Dep) finalTarget() string {
	target := d.Target
	if target == "" {
		target = d.Name
	}
	return target
}

func (d *Dep) recursiveStr() string {
	recursive := "/..."
	if !d.Recursive {
		recursive = ""
	}
	return recursive
}

// CommitBranchTag the thing to checkout
func (d *Dep) CommitBranchTag() string {
	v := ""
	if d.Branch != "" {
		v = d.Branch
	}
	if d.Tag != "" {
		v = d.Tag
	}
	if d.Commit != "" {
		v = d.Commit
	}
	return v
}

func (d *Dep) detectVcsFolder() (string, error) {
	var path string
	path, err := gbutils.FindInAncestorPath(d.pkgVendorFolder(), ".git")
	if err != nil {
		path, err = gbutils.FindInAncestorPath(d.pkgVendorFolder(), ".hg")
		if err != nil {
			path, err = gbutils.FindInAncestorPath(d.pkgVendorFolder(), ".bzr")
			if err != nil {
				return "", fmt.Errorf("not VCS folder found for %s", d.pkgVendorFolder())
			}
		}
	}
	replaceStr := filepath.Join(d.vendorFolder(), "src")
	path = strings.Replace(path, replaceStr, "", 1)
	d.VcsFolder = path
	fmt.Printf("GIT FOLDER: %s\n", path)
	return path, nil
}

func plainRunCmd(dir string, argStr string) error {
	return runCmd(dir, strings.Split(argStr, " "), []string{})
}

// http://craigwickesser.com/2015/02/golang-cmd-with-custom-environment/
func runCmd(dir string, args []string, cmdEnv []string) error {
	env := os.Environ()
	for _, str := range cmdEnv {
		env = append(env, str)
	}
	fmt.Println(fmt.Sprintf("in %s running", dir))
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
