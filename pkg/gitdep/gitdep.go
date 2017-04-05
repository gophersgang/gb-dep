package gitdep

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

// Dep is a (GIT) dependency
type Dep struct {
	packagefile.Pkg
	RootFolder string // the root folder
}

// Run knows what to do
func (d *Dep) Run() error {
	d.ensure()
	fmt.Println(d.Name)
	fmt.Println(d.CommitBranchTag())
	d.ensureCacheFresh()
	d.ensureCheckout()

	return nil
}

// Removevcs removes .GIT folder in vendor,
// run before commiting
func (d *Dep) Removevcs() error {
	return os.RemoveAll(d.checkoutGitFolder())
}

func (d *Dep) checkoutGitFolder() string {
	return filepath.Join(d.pkgCheckoutFolder(), ".git")
}

// Copyvcs creates .GIT folder in vendor, for installation and such
// run before commiting
func (d *Dep) Copyvcs() error {
	if d.vcsFolderCopied() {
		return nil
	}
	copyCmd := fmt.Sprintf("cp -r %s/.git %s/.git", d.pkgCachedFolder(), d.pkgCheckoutFolder())
	return runCmd(d.vendorFolder(), strings.Split(copyCmd, " "), []string{})
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

func (d *Dep) cloneArgs() []string {
	args := []string{"git", "clone"}
	args = append(args, "https://"+d.Name)
	args = append(args, d.pkgCachedFolder())
	return args
}

func (d *Dep) fetchArgs() []string {
	args := []string{"git", "fetch"}
	return args
}

func (d *Dep) clone() error {
	runCmd(d.cacheFolder(), d.cloneArgs(), []string{})
	return nil
}

func (d *Dep) update() error {
	runCmd(d.pkgCachedFolder(), d.fetchArgs(), []string{})
	return nil
}

func (d *Dep) alreadyCloned() bool {
	return fileExist(d.pkgCachedFolder())
}

func (d *Dep) vcsFolderCopied() bool {
	return fileExist(d.checkoutGitFolder())
}

func fileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (d *Dep) ensureCacheFresh() error {
	if d.alreadyCloned() {
		d.update()
	} else {
		d.clone()
	}
	return nil
}

func (d *Dep) ensureCheckout() error {
	d.Copyvcs()
	checkoutCmd := fmt.Sprintf("git reset --hard %s", d.CommitBranchTag())
	runCmd(d.pkgCheckoutFolder(), strings.Split(checkoutCmd, " "), []string{})
	runCmd(d.pkgCheckoutFolder(), []string{"go", "get", "-u", d.Name}, []string{})
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

func (d *Dep) ensure() error {
	os.MkdirAll(d.cacheFolder(), 0777)
	os.MkdirAll(d.vendorFolder(), 0777)
	os.MkdirAll(d.pkgCheckoutFolder(), 0777)
	return nil
}

// http://craigwickesser.com/2015/02/golang-cmd-with-custom-environment/
func runCmd(dir string, args []string, cmdEnv []string) error {
	env := os.Environ()
	for _, str := range cmdEnv {
		env = append(env, str)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = env
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
