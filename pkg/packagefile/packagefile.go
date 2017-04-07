package packagefile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gophersgang/gb-dep/pkg/gbutils"
	hjson "github.com/hjson/hjson-go"
)

var (
	pkgFile = "package.hjson"
)

// PackageFile represent a packagefile
type PackageFile struct {
	Packages []Pkg `json:"packages"`
}

// Pkg represents a Golang package
type Pkg struct {
	Name      string   `json:"name"`                // name of the package
	Group     []string `json:"group,omitempty"`     // in which groups should this package be installed?
	Goos      []string `json:"goos,omitempty"`      // what Go OS is supported
	Insecure  bool     `json:"insecure,omitempty"`  // use insecure protocol for downloading
	Recursive bool     `json:"recursive,omitempty"` // should we also fetch the vendored packages for a package?
	Command   string   `json:"command,omitempty"`   // special command to executed for downloading the packag
	Private   bool     `json:"private,omitempty"`   // is this package private?
	Skipdep   bool     `json:"skipdep,omitempty"`   // shall we ignore dependencies?
	Target    string   `json:"target,omitempty"`    // folder to install package to
	// current status, commit > tag > branch
	Branch string `json:"branch,omitempty"`
	Tag    string `json:"tag,omitempty"`
	Commit string `json:"commit,omitempty"`
}

var allowedFields = []string{
	"name",
	"group",
	"goos",
	"insecure",
	"recursive",
	"command",
	"private",
	"skipdep",
	"target",
	// current status, commit > tag > branch
	"branch",
	"tag",
	"commit",
}

// Parse will read a file and return Pgk structs
func Parse(path string) ([]Pkg, error) {
	content, err := ioutil.ReadFile(path)
	checkErr("Could not read "+path, err)
	if err != nil {
		return nil, err
	}
	pfile, err := myunmarshal(content)
	checkErr("Could not ummarshal "+path, err)
	if err != nil {
		return nil, err
	}
	return pfile.Packages, nil
}

// FindPackagefile returns the path to package.hjson in the given path
func FindPackagefile(dir string) (string, error) {
	return gbutils.FindInAncestorPath(dir, pkgFile)
}

/*
	- convert hjson content to json
	- map it to PackageFile
*/
func myunmarshal(content []byte) (PackageFile, error) {
	var data interface{}
	hjson.Unmarshal(content, &data)
	jsonRaw, err := json.Marshal(data)
	if err != nil {
		return PackageFile{}, err
	}

	err = validatejson(jsonRaw)
	if err != nil {
		return PackageFile{}, err
	}

	into := PackageFile{}
	err = json.Unmarshal(jsonRaw, &into)
	if err != nil {
		return PackageFile{}, err
	}
	return into, nil
}

// because stdlib json ignores unrecognized fields, this is a quick way to
// check that all the packages include only the allowed keys
func validatejson(jsonRaw []byte) error {
	var dat map[string]interface{}
	if err := json.Unmarshal(jsonRaw, &dat); err != nil {
		panic(err)
	}
	packages := dat["packages"].([]interface{})
	for _, pkg := range packages {
		newPkg := pkg.(map[string]interface{})
		for key := range newPkg {
			if !contains(allowedFields, key) {
				return errors.New("KEY NOT ALLOWED " + key)
			}
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// just a quick verbose error printer
func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
	}
}
