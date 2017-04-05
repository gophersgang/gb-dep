package install

import (
	"flag"
	"fmt"
	"log"

	"os"

	"github.com/gophersgang/gb-dep/pkg/packagefile"
	"github.com/gophersgang/gb-dep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("install", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	fmt.Println("Running install....")
	install(args)
	// runner.Run([]string{"echo", "this"}, runner.Green)
}

func (r *cmd) Usage() string {
	return "install --verbose=true"
}

func install(args []string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}
	file, err := packagefile.FindPackagefile(currDir)
	if err != nil {
		return err
	}
	pkgs, err := packagefile.Parse(file)
	if err != nil {
		return err
	}
	for _, pkg := range pkgs {
		fmt.Println(pkg.Name)
	}
	return nil
}

// func install(args []string) error {
// 	allGoms, err := parseGomfile("Gomfile")
// 	if err != nil {
// 		return err
// 	}

// 	vendor, err := filepath.Abs(vendorFolder)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = os.Stat(vendor)
// 	if err != nil {
// 		err = os.MkdirAll(vendor, 0755)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	err = os.Setenv("GOPATH", vendor)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.Setenv("GOBIN", filepath.Join(vendor, "bin"))
// 	if err != nil {
// 		return err
// 	}

// 	// 1. Filter goms to install
// 	goms := make([]Gom, 0)
// 	for _, gom := range allGoms {
// 		if group, ok := gom.options["group"]; ok {
// 			if !matchEnv(group) {
// 				continue
// 			}
// 		}
// 		if goos, ok := gom.options["goos"]; ok {
// 			if !matchOS(goos) {
// 				continue
// 			}
// 		}
// 		goms = append(goms, gom)
// 	}

// 	// 2. Clone the repositories
// 	for _, gom := range goms {
// 		err = gom.Clone(args)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// 3. Checkout the commit/branch/tag if needed
// 	for _, gom := range goms {
// 		err = gom.Checkout()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// 4. Build and install
// 	for _, gom := range goms {
// 		if skipdep, ok := gom.options["skipdep"].(string); ok {
// 			if skipdep == "true" {
// 				continue
// 			}
// 		}
// 		err = gom.build(args)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
