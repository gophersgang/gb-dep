package packagefile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/gbutils"
	"github.com/gophersgang/gbdep/pkg/structs"

	hjson "github.com/hjson/hjson-go"
)

var (
	pkgFile  = "package.hjson"
	md5Field = "package_md5"
)

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
	"locked_commit",
	"vcs_type",
}

// GimmePackages is a higlevel function, that hides the implementation details of how you get
// packages information. It works on current path
func GimmePackages() ([]structs.Pkg, error) {
	currDir, err := os.Getwd()
	if err != nil {
		return []structs.Pkg{}, err
	}
	file, err := FindPackagefile(currDir)
	if err != nil {
		return []structs.Pkg{}, err
	}
	// root := filepath.Dir(file)
	pkgs, err := Parse(file)
	if err != nil {
		return nil, err
	}
	return pkgs, nil
}

// Parse will read a file and return Pgk structs
func Parse(path string) ([]structs.Pkg, error) {
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

// GenerateLockFile creates a lockfile
func GenerateLockFile(deps []*dep.Dep) error {
	currDir, err := os.Getwd()
	packFile, err := FindPackagefile(currDir)
	lockpath := filepath.Join(filepath.Dir(packFile), "package.lock")

	this := map[string]interface{}{"packages": deps}
	md5, err := gbutils.ComputeMD5(packFile)
	this[md5Field] = md5
	res, err := json.MarshalIndent(this, "", "    ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(lockpath, res, 0777)
	return nil
}

// RootDir returns the folder where package.hjson is found
func RootDir(dir string) (string, error) {
	packFile, err := FindPackagefile(dir)
	if err != nil {
		return "", err
	}
	root := filepath.Dir(packFile)
	return root, nil
}

/*
	- convert hjson content to json
	- map it to PackageFile
*/
func myunmarshal(content []byte) (structs.PackageFile, error) {
	var data interface{}
	hjson.Unmarshal(content, &data)
	jsonRaw, err := json.Marshal(data)
	if err != nil {
		return structs.PackageFile{}, err
	}

	err = validatejson(jsonRaw)
	if err != nil {
		return structs.PackageFile{}, err
	}

	into := structs.PackageFile{}
	err = json.Unmarshal(jsonRaw, &into)
	if err != nil {
		return structs.PackageFile{}, err
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
			if !gbutils.ContainsStr(allowedFields, key) {
				return errors.New("KEY NOT ALLOWED " + key)
			}
		}
	}

	return nil
}

// just a quick verbose error printer
func checkErr(msg string, err error) {
	if err != nil {
		fmt.Println(msg)
		fmt.Println(err)
	}
}
