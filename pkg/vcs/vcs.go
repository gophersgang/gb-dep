package vcs

// alternatives:
// https://github.com/govend/govend/blob/master/deps/vcs/vcs.go
import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type VcsCmd struct {
	checkout     []string
	update       []string
	revision     []string
	revisionMask string
}

var (
	HG = &VcsCmd{
		[]string{"hg", "update"},
		[]string{"hg", "pull"},
		[]string{"hg", "id", "-i"},
		"^(.+)$",
	}
	GIT = &VcsCmd{
		[]string{"git", "checkout", "-q"},
		[]string{"git", "fetch"},
		[]string{"git", "rev-parse", "HEAD"},
		"^(.+)$",
	}
	BZR = &VcsCmd{
		[]string{"bzr", "revert", "-r"},
		[]string{"bzr", "pull"},
		[]string{"bzr", "log", "-r-1", "--line"},
		"^([0-9]+)",
	}

	Versions = map[string]*VcsCmd{
		"hg":  HG,
		"git": GIT,
		"bzr": BZR,
	}
)

func (vcs *VcsCmd) Checkout(p, destination string) error {
	args := append(vcs.checkout, destination)
	return VcsExec(p, args...)
}

func (vcs *VcsCmd) Update(p string) error {
	return VcsExec(p, vcs.update...)
}

func (vcs *VcsCmd) Revision(dir string) (string, error) {
	args := append(vcs.revision)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}
	rev := strings.TrimSpace(string(b))
	if vcs.revisionMask != "" {
		return regexp.MustCompile(vcs.revisionMask).FindString(rev), nil
	}
	return rev, nil
}

func (vcs *VcsCmd) Sync(p, destination string) error {
	err := vcs.Checkout(p, destination)
	if err != nil {
		err = vcs.Update(p)
		if err != nil {
			return err
		}
		err = vcs.Checkout(p, destination)
	}
	return err
}

func VcsExec(dir string, args ...string) error {
	fmt.Println("...VcsExec...")
	fmt.Printf("IN %s\n", dir)
	fmt.Println(args)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
