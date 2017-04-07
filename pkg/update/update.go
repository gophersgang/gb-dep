package update

import (
	"flag"
	"fmt"
	"log"

	"github.com/gophersgang/gb-dep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("update", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	fmt.Println("Running update....")
}

func (r *cmd) Usage() string {
	return "update --verbose=true"
}

// func update() error {
// 	goms, err := parseGomfile("Gomfile")
// 	if err != nil {
// 		return err
// 	}
// 	vendor, err := filepath.Abs(vendorFolder)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.Setenv("GOPATH", vendor)
// 	if err != nil {
// 		return err
// 	}
// 	err = os.Setenv("GOBIN", filepath.Join(vendor, "bin"))
// 	if err != nil {
// 		return err
// 	}

// 	for _, gom := range goms {
// 		err = gom.Update()
// 		if err != nil {
// 			return err
// 		}
// 		vcs, _, p := vcsScan(vendorSrc(vendor), gom.name)
// 		if vcs != nil {
// 			rev, err := vcs.Revision(p)
// 			if err == nil && rev != "" {
// 				gom.options["commit"] = rev
// 			}
// 		}
// 	}

// 	return writeGomfile("Gomfile", goms)
// }
