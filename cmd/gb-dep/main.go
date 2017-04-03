package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gophersgang/gb-dep/install"
	"github.com/gophersgang/gb-dep/subcommands"
	"github.com/gophersgang/gb-dep/update"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	subCmds := subcommands.New(os.Args[0], "A CLI GB-dep")

	subCmds.Register("update", "Update dependencies", update.New())
	subCmds.Register("install", "Install dependencies", install.New())

	quiet := flag.Bool("quiet", false, "Silence output")
	flag.Parse()
	logger := log.New(ioutil.Discard, "", 0)
	log.SetOutput(ioutil.Discard)

	if !*quiet {
		logger = log.New(os.Stderr, "", 0)
	}

	subCmds.Run(flag.Args(), logger)
}
