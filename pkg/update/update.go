package update

import (
	"flag"
	"fmt"
	"log"

	"github.com/gophersgang/gbdep/pkg/subcommands"
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
