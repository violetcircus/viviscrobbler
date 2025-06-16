package setup

import (
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type SessionResponse struct {
	Session Session `json:"session"`
}

type Session struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func Setup() {
	createConfig()
}

// create the files in .config needed by the program if they don't exist
func createConfig() {
	files := []string{"config.toml", ".lastfm_session", "logFile.tsv", "mapFile.tsv"}
	fmt.Println("hi")

	for _, file := range files {
		if _, err := os.Stat(configreader.GetConfigDir() + file); err == nil {
			continue // skip file if it exists already
		} else if errors.Is(err, os.ErrNotExist) {
			// make it if it doesn't
			configFile, err := os.Create(configreader.GetConfigDir() + file)
			if err != nil {
				log.Fatal(err)
			}
			defer configFile.Close()

			// do different things per file
			switch file {
			// if there's no saved last fm session, make one
			case ".lastfm_session":
				requestAuth()
			// write default config to file
			case "config.toml":
				WriteConfig()
			}
		} else {
			// handle error case that isnt the file not existing
			log.Println("file may or may not exist. Weird")
		}
	}
}

// takes a map of url parameters and creates the signature.
func SignSignature(parameters map[string]string) string {
	secrets := secret.GetSecrets()

	// get parameter keys and sort them
	keys := make([]string, len(parameters))
	i := 0
	for k := range parameters {
		keys[i] = k
		i++ // using a count int instead of an append() is more memory efficient #genius
	}
	sort.Strings(keys)

	// build signature base string
	var sb strings.Builder
	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteString(parameters[key])
	}
	sb.WriteString(secrets.Secret)

	// calculate md5 hash
	h := md5.New()
	io.WriteString(h, sb.String())
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

// gets an auth token from Last.FM
func GetToken() string {
	parameters := map[string]string{
		"api_key": secret.GetSecrets().ApiKey,
	}
	signature := SignSignature(parameters)
	baseUrl := "https://ws.audioscrobbler.com/2.0/"
	urlParams := fmt.Sprintf("?method=auth.gettoken&api_key=%v&api_sig=%v&format=json", parameters["api_key"], signature)

	resp, err := http.Get(baseUrl + urlParams)
	if err != nil {
		log.Println("get error")
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read error")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var result TokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("json error")
		log.Fatal(err)
	}
	token := strings.TrimSpace(result.Token)

	return token
}

// request authorisation from last fm
func requestAuth() {
	token := GetToken()
	api_key := secret.GetSecrets().ApiKey
	fmt.Println("Click the link to authorise viviscrobbler with Last.FM, then press enter:")
	fmt.Printf("http://www.last.fm/api/auth/?api_key=%v&token=%v", api_key, token)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	getSession(token)
}

// get last fm session id and username
func getSession(token string) {
	apiKey := secret.GetSecrets().ApiKey
	parameters := map[string]string{
		"api_key": apiKey,
		"method":  "auth.getSession",
		"token":   token,
	}
	signature := SignSignature(parameters)

	baseUrl := "https://ws.audioscrobbler.com/2.0/"
	urlParams := fmt.Sprintf("?method=auth.getSession&api_key=%v&token=%v&api_sig=%v&format=json", apiKey, token, signature)

	resp, err := http.Get(baseUrl + urlParams)
	if err != nil {
		log.Println("get error")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("read error")
		log.Fatal(err)
	}
	var result SessionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("json error")
		log.Fatal(err)
	}
	writeSession(result.Session)
}

// write the session to a file that can be read from later
func writeSession(session Session) {
	f, err := os.Create(configreader.GetConfigDir() + ".lastfm_session")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	row := []string{session.Name, session.Key}
	if err := w.Write(row); err != nil {
		log.Fatal("error writing to file")
	}
	w.Flush()
}

// writes the default config
func WriteConfig() {
	config := configreader.Config{
		SingleArtist:      true,
		ApiCheck:          true,
		Regex:             "",
		ScrobbleThreshold: 50.0,
		ApiKey:            "",
		Secret:            "",
	}
	f, err := os.Create(configreader.GetConfigDir() + "config.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	t := toml.NewEncoder(f)
	if err := t.Encode(config); err != nil {
		log.Fatal(err)
	}
	err = f.Sync()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("config written!")
}
