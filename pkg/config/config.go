package config

import (
	"path/filepath"

	"log"
	"os"

	"github.com/comail/colog"
)

var (
	// Config present our assumptions
	Config = &config{
		VendorFolder: "vendor/src",
	}
)

type config struct {
	VendorFolder  string
	Logger        *log.Logger
	LoggerBackend *colog.CoLog
}

func (cfg *config) AbsVendorFolder() (string, error) {
	vendor, err := filepath.Abs(cfg.VendorFolder)
	if err != nil {
		return "", err
	}
	return vendor, nil
}

func (cfg *config) SetDebugMode() {
	cfg.LoggerBackend.SetMinLevel(colog.LDebug)
}

func init() {
	cl := colog.NewCoLog(os.Stdout, "install ", log.LstdFlags)
	cl.SetMinLevel(colog.LInfo)
	cl.SetDefaultLevel(colog.LDebug)
	//cl.SetDefaultLevel(colog.LWarning)
	logger := cl.NewLogger()
	Config.Logger = logger
	Config.LoggerBackend = cl
}
