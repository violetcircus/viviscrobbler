package configreader

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
)

type Config struct {
	ServerAddress     string
	ServerPort        string
	SingleArtist      bool
	ApiCheck          bool
	Regex             string
	ScrobbleThreshold float64
	ApiKey            string
	Secret            string
}

// get config directory
func GetConfigDir() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	return configDir + "/vvscrobbler/"
}

// read config
func ReadConfig() Config {
	var configfile = GetConfigDir() + "config.toml"
	x, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing:", configfile)
	}
	_ = x

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}
