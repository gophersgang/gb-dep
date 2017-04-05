package packagefile_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/gophersgang/gb-dep/pkg/packagefile"
)

func execCommands(dir string, args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TestFindFile(t *testing.T) {
	os.MkdirAll("/tmp/this/is/a/very/long/path", 0777)
	ioutil.WriteFile("/tmp/this/is/a/package.hjson", []byte{0}, 0777)

	res, err := packagefile.FindPackagefile("/tmp/this/is/a/very/long/path")
	if err != nil {
		t.Errorf("Lookup should not fail...")
	}
	expected := "/tmp/this/is/a/package.hjson"
	if res != expected {
		t.Errorf("Could not find package file at %s, was at %s", expected, res)
	}
}

func TestNotFound(t *testing.T) {
	_, err := packagefile.FindPackagefile("/tmp/does/not/exist")
	if err == nil {
		t.Errorf("Lookup should fail")
	}
}
