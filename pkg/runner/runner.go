package runner

import (
	"bytes"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/daviddengcn/go-colortext"
)

type Color int

const (
	None  Color = Color(ct.None)
	Red   Color = Color(ct.Red)
	Blue  Color = Color(ct.Blue)
	Green Color = Color(ct.Green)
)

var (
	vendorFolder = "src/vendor"
)

func handleSignal() {
	sc := make(chan os.Signal, 10)
	signal.Notify(sc, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		<-sc
		ct.ResetColor()
		os.Exit(0)
	}()
}

// will prepare
func Ready() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	vendor, err := filepath.Abs(vendorFolder)
	if err != nil {
		return err
	}

	for {
		file := filepath.Join(dir, "package.hjson")
		if isFile(file) {
			vendor = filepath.Join(dir, vendorFolder)
			break
		}
		next := filepath.Clean(filepath.Join(dir, ".."))
		if next == dir {
			dir = ""
			break
		}
		dir = next
	}

	binPath := os.Getenv("PATH") +
		string(filepath.ListSeparator) +
		filepath.Join(vendor, "bin")
	err = os.Setenv("PATH", binPath)
	if err != nil {
		return err
	}

	var paths []string
	if dir == "" {
		paths = []string{vendor, os.Getenv("GOPATH")}
	} else {
		paths = []string{vendor, dir, os.Getenv("GOPATH")}
	}
	vendor = strings.Join(paths, string(filepath.ListSeparator))
	err = os.Setenv("GOPATH", vendor)
	if err != nil {
		return err
	}

	return nil
}

var stdout = os.Stdout
var stderr = os.Stderr
var stdin = os.Stdin

// Run runs a given commands with given color
func Run(args []string, c Color) error {
	if err := Ready(); err != nil {
		return err
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	ct.ChangeColor(ct.Color(c), true, ct.None, false)
	err := cmd.Run()
	ct.ResetColor()
	return err
}

func List(dir string) ([]string, error) {
	cmd := exec.Command("go", "list", "./...")
	cmd.Dir = dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return strings.Split(stdout.String(), "\n"), nil
}

func isFile(p string) bool {
	if fi, err := os.Stat(filepath.Join(p)); err == nil && !fi.IsDir() {
		return true
	}
	return false
}
