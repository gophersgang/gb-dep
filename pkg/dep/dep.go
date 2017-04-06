package dep

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gophersgang/gb-dep/pkg/config"
	"github.com/gophersgang/gb-dep/pkg/packagefile"
)

var (
	cfg = config.Config
)

// Dep is a package
type Dep struct {
	packagefile.Pkg
	RootFolder string // the root folder
}

// Run knows what to do
func (d *Dep) Run() error {
	fmt.Println("Installing " + d.Name)
	d.ensureBasefolders()
	d.ensureCheckout()
	return nil
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

func fileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (d *Dep) ensureCheckout() error {
	os.Setenv("GOBIN", filepath.Join(d.vendorFolder(), "bin")) // this is important to get binaries
	os.Setenv("GOPATH", filepath.Join(d.vendorFolder()))       // so it works reliably

	myRun := func(dir string, argStr string) {
		runCmd(dir, strings.Split(argStr, " "), []string{})
	}
	myRun(d.vendorFolder(), fmt.Sprintf("go get -u %s", d.Name))
	myRun(d.pkgCheckoutFolder(), "go get -u ./...")
	myRun(d.pkgCheckoutFolder(), "go install ./...")
	return nil
}

func (d *Dep) pkgCachedFolder() string {
	return filepath.Join(d.cacheFolder(), d.Name)
}
func (d *Dep) pkgCheckoutFolder() string {
	return filepath.Join(d.vendorFolder(), "src", d.Name)
}

func (d *Dep) cacheFolder() string {
	return d.RootFolder + "/" + "vendor.cache"
}

func (d *Dep) vendorFolder() string {
	return d.RootFolder + "/" + "vendor"
}

func (d *Dep) ensureBasefolders() error {
	os.MkdirAll(d.cacheFolder(), 0777)
	os.MkdirAll(d.vendorFolder(), 0777)
	return nil
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
