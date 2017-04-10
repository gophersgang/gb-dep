package cmdcleanvcs

import (
	"flag"
	"log"
	"os"

	"github.com/gophersgang/gbdep/pkg/cmdcommon"
	"github.com/gophersgang/gbdep/pkg/config"
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
	r.fs = flag.NewFlagSet("cleanvcs", flag.ExitOnError)
	r.fs.BoolVar(&r.verbose, "verbose", false, "Noisy output")

	return &r
}

func (r *cmd) Run(args []string, log *log.Logger) {
	r.fs.Parse(args)
	cfg.LoggerBackend.SetPrefix("cleanvcs ")
	if r.verbose {
		cfg.SetDebugMode()
	}
	cfg.Logger.Println("Removing VSC folders from vendor...")
	realcmd(args)
}

func (r *cmd) Usage() string {
	return "cleanvcs --verbose=true"
}

func realcmd(args []string) error {
	folders, err := cmdcommon.AllVendoredVCSFolders()
	if err != nil {
		cfg.Logger.Fatal(err)
	}
	for _, folder := range folders {
		cfg.Logger.Printf("removing %s", folder)
		err := os.RemoveAll(folder)
		if err != nil {
			return err
		}
	}
	// folders
	// cmdcommon.RunConcurrently(deps, 5, func(d *dep.Dep) {
	// 	d.CleanVCS()
	// })
	return nil
}
