package dep

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gophersgang/gb-dep/pkg/packagefile"
	. "github.com/gophersgang/gb-dep/pkg/runner"
	"github.com/gophersgang/gb-dep/pkg/vcs"
)

var (
	vendorFolder = "src/vendor"
)

type Dep struct {
	packagefile.Pkg
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

func (d *Dep) commitBranchTag() string {
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

func (d *Dep) Update() error {
	cmdArgs := []string{"go", "get", "-u"}
	if d.Insecure {
		cmdArgs = append(cmdArgs, "-insecure")
	}

	recursive := d.recursiveStr()
	cmdArgs = append(cmdArgs, d.Name+recursive)

	fmt.Printf("updating %s\n", d.Name)
	return Run(cmdArgs, Green)
}

func (d *Dep) Clone(args []string) error {
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	if d.Command != "" {
		target := d.finalTarget()

		srcdir := filepath.Join(vendor, "src", target)
		if err := os.MkdirAll(srcdir, 0755); err != nil {
			return err
		}

		customCmd := strings.Split(d.Command, " ")
		customCmd = append(customCmd, srcdir)

		fmt.Printf("fetching %s (%v)\n", d.Name, customCmd)
		err = Run(customCmd, Blue)
		if err != nil {
			return err
		}
	}

	if d.Skipdep {
		return nil
	}

	cmdArgs := []string{"go", "get", "-d"}
	if d.Insecure {
		cmdArgs = append(cmdArgs, "-insecure")
	}

	recursive := d.recursiveStr()
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, d.Name+recursive)

	fmt.Printf("downloading %s\n", d.Name)
	return Run(cmdArgs, Blue)
}

func (d *Dep) Checkout() error {
	commit_or_branch_or_tag := d.commitBranchTag()

	if commit_or_branch_or_tag == "" {
		return nil
	}
	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}
	p := filepath.Join(vendor, "src")
	target := d.finalTarget()

	for _, elem := range strings.Split(target, "/") {
		var dvcs *vcs.VcsCmd
		p = filepath.Join(p, elem)
		if isDir(filepath.Join(p, ".git")) {
			dvcs = vcs.GIT
		} else if isDir(filepath.Join(p, ".hg")) {
			dvcs = vcs.HG
		} else if isDir(filepath.Join(p, ".bzr")) {
			dvcs = vcs.BZR
		}
		if dvcs != nil {
			p = filepath.Join(vendor, "src", target)
			return dvcs.Sync(p, commit_or_branch_or_tag)
		}
	}
	fmt.Printf("Warning: don't know how to checkout for %v\n", d.Name)
	return errors.New("gom currently support git/hg/bzr for specifying tag/branch/commit")
}

func hasGoSource(p string) bool {
	dir, err := os.Open(p)
	if err != nil {
		return false
	}
	defer dir.Close()
	fis, err := dir.Readdir(-1)
	if err != nil {
		return false
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			return true
		}
	}
	return false
}

func (d *Dep) build(args []string) (err error) {
	var vendor string
	vendor, err = filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}

	installCmd := []string{"go", "get"}
	hasPkg := false
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			arg = path.Join(arg, "...")
			hasPkg = true
		}
		installCmd = append(installCmd, arg)
	}

	target := d.finalTarget()
	pkgPath := filepath.Join(vendor, "src", target)

	if hasPkg {
		return vcs.VcsExec(pkgPath, installCmd...)
	}

	pkgs, err := List(pkgPath)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if isIgnorePackage(pkg) {
			continue
		}
		pkgPath = filepath.Join(vendor, "src", pkg)
		if !hasGoSource(pkgPath) {
			continue
		}
		err := vcs.VcsExec(pkgPath, installCmd...)
		if err != nil {
			return err
		}
	}
	return nil
}

func isFile(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && !fi.IsDir() {
		return true
	}
	return false
}

func isDir(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && fi.IsDir() {
		return true
	}
	return false
}

func isIgnorePackage(pkg string) bool {
	if pkg == "" {
		return true
	}
	paths := strings.Split(pkg, "/")
	for _, path := range paths {
		if path == "examples" {
			return true
		}
		if strings.HasPrefix(path, "_") {
			return true
		}
	}
	return false
}

func readdirnames(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	return list, nil
}

func parseInstallFlags(args []string) (opts map[string]string, retargs []string) {
	opts = make(map[string]string)
	re := regexp.MustCompile(`^--([a-z][a-z_]*)(=\S*)?`)
	for _, arg := range args {
		ss := re.FindAllStringSubmatch(arg, -1)
		if len(ss) > 0 {
			opts[ss[0][1]] = opts[ss[0][2]]
		} else {
			retargs = append(retargs, arg)
		}
	}
	return
}

func hasSaveOpts(opts map[string]string) bool {
	if _, ok := opts["save"]; ok {
		return true
	}
	if _, ok := opts["save-dev"]; ok {
		return true
	}
	return false
}
