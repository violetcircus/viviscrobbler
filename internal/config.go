package internal

import (
	"log"
	"github.com/BurntSushi/toml" // switch to yaml! need better key:value support
	// "gopkg.in/yaml.v3" // use this prob  
	"os"
)

type Config struct {
	Service string
	SingleArtist bool
	SanityCheck bool
	Regex string
}

func ReadConfig() Config {
	var configfile = "/home/violet/.config/vvscrob/config.toml"
	_, err := os.Stat(configfile)
	if err!= nil {
		log.Fatal("Config file is missing:", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
