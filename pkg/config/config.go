package config

import "path/filepath"

var (
	// Config present our assumptions
	Config = &config{
		VendorFolder: "vendor/src",
	}
)

type config struct {
	VendorFolder string
}

func (cfg *config) AbsVendorFolder() (string, error) {
	vendor, err := filepath.Abs(cfg.VendorFolder)
	if err != nil {
		return "", err
	}
	return vendor, nil
}
