package update

import (
	"flag"
	"fmt"
	"log"

	"github.com/gophersgang/gb-dep/pkg/subcommands"
)

type runner struct {
	fs      *flag.FlagSet
	verbose bool
}

func New() subcommands.Command {
	r := runner{}
	r.fs = flag.NewFlagSet("update", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *runner) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	fmt.Println("Running update....")
}

func (r *runner) Usage() string {
	return "update --verbose=true"
}
