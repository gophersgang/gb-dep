package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gophersgang/gbdep/pkg/cmdbuildbins"
	"github.com/gophersgang/gbdep/pkg/cmdcleanvcs"
	"github.com/gophersgang/gbdep/pkg/cmdinstall"
	"github.com/gophersgang/gbdep/pkg/cmdupdate"
	"github.com/gophersgang/gbdep/pkg/subcommands"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	subCmds := subcommands.New(os.Args[0], "A CLI gbdep")

	subCmds.Register("update", "Update dependencies", cmdupdate.New())
	subCmds.Register("install", "Install dependencies", cmdinstall.New())
	subCmds.Register("cleanvcs", "Removes .git and such from vendor folder", cmdcleanvcs.New())
	subCmds.Register("buildbins", "Builds vendor binaries / libraries", cmdbuildbins.New())

	quiet := flag.Bool("quiet", false, "Silence output")
	flag.Parse()
	logger := log.New(ioutil.Discard, "", 0)
	log.SetOutput(ioutil.Discard)

	if !*quiet {
		logger = log.New(os.Stderr, "", 0)
	}

	subCmds.Run(flag.Args(), logger)
}
