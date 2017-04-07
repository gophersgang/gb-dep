package install

import (
	"flag"
	"log"

	"github.com/gophersgang/gbdep/pkg/cmdcommon"
	"github.com/gophersgang/gbdep/pkg/config"
	"github.com/gophersgang/gbdep/pkg/dep"
	"github.com/gophersgang/gbdep/pkg/packagefile"
	"github.com/gophersgang/gbdep/pkg/subcommands"
)

var (
	cfg = config.Config
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

func (r *cmd) Usage() string {
	return "install --verbose=true"
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	cfg.LoggerBackend.SetPrefix("install ")
	if r.verbose {
		cfg.SetDebugMode()
	}
	cfg.Logger.Println("Installing...")
	realcmd(args)
}

func realcmd(args []string) error {
	deps := cmdcommon.CurrentDeps()
	cmdcommon.RunConcurrently(deps, 5, func(d *dep.Dep) {
		d.Run()
	})
	packagefile.GenerateLockFile(deps)
	return nil
}
