package configreader

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	Service           string
	SingleArtist      bool
	SanityCheck       bool
	ApiCheck          bool
	Regex             string
	ScrobbleThreshold float64
}

func ReadConfig() Config {
	var configfile = "/home/violet/.config/vvscrob/config.toml"
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing:", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
