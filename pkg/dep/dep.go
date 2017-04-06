package dep

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gophersgang/gb-dep/pkg/config"
	"github.com/gophersgang/gb-dep/pkg/gbutils"
	"github.com/gophersgang/gb-dep/pkg/packagefile"
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
	if d.cacheExists() {
		return nil
	}
	copyCmd := fmt.Sprintf("cp -r %s/.git %s/.git", d.pkgVendorFolder(), d.pkgCachedFolder())
	return plainRunCmd(d.vendorFolder(), copyCmd)
}

// func (d *Dep) vcsFolderCached() bool {
// 	return fileExist(d.pkgCachedFolder())
// }

func (d *Dep) cacheExists() bool {
	return gbutils.PathExist(filepath.Join(d.pkgCachedFolder(), ".git"))
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

func (d *Dep) detectVcsFolder() string {

	return ""
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
