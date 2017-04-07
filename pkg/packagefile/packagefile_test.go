package packagefile_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gophersgang/gbdep/pkg/packagefile"
)

func tempGomfile(content string) (string, error) {
	f, err := ioutil.TempFile("", "gom")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	name := f.Name()
	return name, nil
}

func TestParse(t *testing.T) {
	filename, err := tempGomfile(`
		packages: [
			// default packages

			// dev packages
			{ name: "github.com/mattn/gover", goos: [ "windows", "linux", "darwin" ], commit: "x8948594854" , group: ["production"]}

			// test packages
			{ name: "github.com/mattn/gom", goos: [ "windows", "linux", "darwin" ], commit: "x8948594854", group: ["test"] } // blabkabla
		]
`)
	if err != nil {
		t.Fatal(err)
	}
	packagef, err := packagefile.Parse(filename)
	pkgs := packagef.Packages
	fmt.Println(pkgs)
	expected := "github.com/mattn/gover"
	if pkgs[0].Name != expected {
		t.Fatalf("Expected %v, but was %v:", expected, pkgs[0].Name)
	}
}

func TestFieldValidation(t *testing.T) {
	filename, _ := tempGomfile(`{
		packages: [
			// default packages

			// dev packages
			{ namesss: "github.com/mattn/gover", goos: [ "windows", "linux", "darwin" ], commit: "x8948594854" , group: ["production"]}

			// test packages
			{ name: "github.com/mattn/gom", goos: [ "windows", "linux", "darwin" ], commit: "x8948594854", group: ["test"] } // blabkabla
		]
	}
`)
	_, err := packagefile.Parse(filename)
	if err == nil {
		t.Fatal("expected error on parsing here")
	}
}

func TestRealFile(t *testing.T) {
	pkgpackagef, _ := packagefile.Parse("assets/package.hjson")
	pkgs := pkgpackagef.Packages
	expected := "github.com/mattn/gover"
	if pkgs[0].Name != expected {
		t.Fatalf("Expected %v, but was %v:", expected, pkgs[0].Name)
	}
}
