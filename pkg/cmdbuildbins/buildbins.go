package cmdbuildbins

import (
	"flag"
	"log"

	"github.com/gophersgang/gbdep/pkg/cmdcommon"
	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/subcommands"
)

type cmd struct {
	fs      *flag.FlagSet
	verbose bool
}

var (
	cfg = config.Config
)

func New() subcommands.Command {
	r := cmd{}
	r.fs = flag.NewFlagSet("buildbins", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	cfg.LoggerBackend.SetPrefix("buildbins ")
	if r.verbose {
		cfg.SetDebugMode()
	}
	cfg.Logger.Println("Build all the vendor binaries / libraries")
	realcmd(args)
}

func (r *cmd) Usage() string {
	return "buildbins --verbose=true"
}

func realcmd(args []string) error {
	deps := cmdcommon.CurrentDeps()
	cmdcommon.RunConcurrently(deps, 5, func(d *dep.Dep) {
		d.CleanVCS()
	})
	packagefile.GenerateLockFile(deps)
	return nil
}
