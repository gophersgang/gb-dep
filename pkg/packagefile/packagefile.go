package packagefile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/gbutils"
	"github.com/gophersgang/gbdep/pkg/structs"

	hjson "github.com/hjson/hjson-go"
)

var (
	pkgFile  = "package.hjson"
	md5Field = "package_md5"
	cfg      = config.Config
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

// GimmePackagefile is a high-level function, that hides the implementation details of how you get
// packages information. It works on current path
func GimmePackagefile() (structs.PackageFile, error) {
	emptyPackagefile := structs.PackageFile{}
	currDir, err := os.Getwd()
	if err != nil {
		return emptyPackagefile, err
	}

	file, err := FindPackagefile(currDir)
	if err != nil {
		return emptyPackagefile, err
	}

	// data from the lockfile
	if isLockfileUptodate(currDir) {
		lockfile, err := FindPackageLock(currDir)
		res, err := Parse(lockfile)
		if err != nil {
			return emptyPackagefile, err
		}
		return res, nil
	}

	// data from normal packagefile
	res, err := Parse(file)
	if err != nil {
		return emptyPackagefile, err
	}
	return res, nil
}

func isLockfileUptodate(dir string) bool {
	file, err := FindPackagefile(dir)
	if err != nil {
		log.Fatal("NOPE")
	}
	lockfile, err := FindPackageLock(dir)
	if err != nil {
		return false
	}
	res, err := Parse(lockfile)
	currentMD5, err := gbutils.ComputeMD5(file)

	if res.PackageMD5 != currentMD5 {
		cfg.Logger.Print("debug: *** STALE lockfile ****")
		return false
	}
	cfg.Logger.Print("debug: *** UP-TO-DATE lockfile ****")
	return true
}

// Parse will read a file and return PackageFile struct
func Parse(path string) (structs.PackageFile, error) {
	content, err := ioutil.ReadFile(path)
	checkErr("Could not read "+path, err)
	if err != nil {
		return structs.PackageFile{}, err
	}
	pfile, err := myunmarshal(content)
	checkErr("Could not ummarshal "+path, err)
	if err != nil {
		return structs.PackageFile{}, err
	}
	return pfile, nil
}

// FindPackagefile returns the path to package.hjson in the given path
func FindPackagefile(dir string) (string, error) {
	return gbutils.FindInAncestorPath(dir, pkgFile)
}

// FindPackageLock returns the path to package.lock in the given path
func FindPackageLock(dir string) (string, error) {
	packFile, err := FindPackagefile(dir)
	if err != nil {
		return "", err
	}
	lockpath := filepath.Join(filepath.Dir(packFile), "package.lock")
	if gbutils.IsFile(lockpath) {
		return lockpath, nil
	}
	return "", fmt.Errorf("No lockfile found for %s", dir)
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
