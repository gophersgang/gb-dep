package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gophersgang/gbdep/pkg/buildbins"
	"github.com/gophersgang/gbdep/pkg/cleanvcs"
	"github.com/gophersgang/gbdep/pkg/install"
	"github.com/gophersgang/gbdep/pkg/subcommands"
	"github.com/gophersgang/gbdep/pkg/update"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	subCmds := subcommands.New(os.Args[0], "A CLI gbdep")

	subCmds.Register("update", "Update dependencies", update.New())
	subCmds.Register("install", "Install dependencies", install.New())
	subCmds.Register("cleanvcs", "Removes .git and such from vendor folder", cleanvcs.New())
	subCmds.Register("buildbins", "Builds vendor binaries / libraries", buildbins.New())

	quiet := flag.Bool("quiet", false, "Silence output")
	flag.Parse()
	logger := log.New(ioutil.Discard, "", 0)
	log.SetOutput(ioutil.Discard)

	if !*quiet {
		logger = log.New(os.Stderr, "", 0)
	}

	subCmds.Run(flag.Args(), logger)
}
