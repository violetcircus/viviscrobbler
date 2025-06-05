package secret

import (
	"embed"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"log"
)

type Secrets struct {
	ApiKey string
	Secret string
}

// embed all files in the secrets folder
// go:embed *
var content embed.FS

// embed api key and secret at compile time, these are in the gitignore so to build this locally
// you will need to add your own api key and secret (acquired from registering an app on LastFM)
// to the config file.
func GetSecrets() Secrets {
	s := Secrets{}
	secretFile, err := content.ReadFile("secret")
	secret := ""
	if err != nil {
		log.Println("Error reading file:", err)
		if configreader.ReadConfig().Secret != "" {
			secret = configreader.ReadConfig().Secret
		} else {
			log.Fatal("your config file needs an Api Key and Secret, read the readme.")
		}
	} else {
		secret = string(secretFile)
	}
	s.Secret = secret

	apiKeyFile, err := content.ReadFile("apiKey")
	apiKey := ""
	if err != nil {
		log.Println("Error reading file:", err)
		if configreader.ReadConfig().ApiKey != "" {
			apiKey = configreader.ReadConfig().ApiKey
		} else {
			log.Fatal("your config file needs an Api Key and Secret, read the readme.")
		}
	} else {
		apiKey = string(apiKeyFile)
	}
	s.ApiKey = apiKey
	return s
}
