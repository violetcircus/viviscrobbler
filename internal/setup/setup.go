package setup

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
)

type TokenResponse struct {
	Token string `json:"token"`
}
type SessionResponse struct {
	Name string `xml:"name"`
	Key  string `xml:"key"`
}

func Setup() {
	createConfig()
}

// create the files in .config needed by the program if they don't exist
func createConfig() {
	files := []string{"config.toml", ".lastfm_session", "logFile.tsv"}

	for _, file := range files {
		if _, err := os.Stat(configreader.ConfigLocation + file); err == nil {
			continue // skip file if it exists already
		} else if errors.Is(err, os.ErrNotExist) {
			// make it if it doesn't
			configFile, err := os.Create(configreader.ConfigLocation + file)
			if err != nil {
				log.Fatal(err)
			}
			defer configFile.Close()

			// if there's no saved last fm session, make one
			if file == ".lastfm_session" {
				requestAuth()
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
		i++
	}
	sort.Strings(keys)

	h := md5.New()
	for _, key := range keys {
		io.WriteString(h, key)
		io.WriteString(h, parameters[key])
	}
	io.WriteString(h, secrets.Secret)
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

// gets an auth token from Last.FM
func GetToken() string {
	parameters := make(map[string]string)
	parameters["apiKey"] = secret.GetSecrets().ApiKey
	signature := SignSignature(parameters)
	baseUrl := "https://ws.audioscrobbler.com/2.0/"
	urlParams := fmt.Sprintf("?method=auth.gettoken&api_key=%v&%v&format=json", parameters["apiKey"], signature)

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
	var result TokenResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("json error")
		log.Fatal(err)
	}
	token := result.Token

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
	parameters := make(map[string]string)
	parameters["token"] = token
	parameters["apiKey"] = apiKey
	signature := SignSignature(parameters)
	baseUrl := "http://www.last.fm/api/auth/"
	urlParams := fmt.Sprintf("?method=auth.getsession&api_key=%v&token=%v&%v", apiKey, token, signature)

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
	fmt.Print(body)
	var result SessionResponse
	if err := xml.Unmarshal(body, &result); err != nil {
		log.Println("xml error")
		log.Fatal(err)
	}
	writeSession(result)
}

func writeSession(result SessionResponse) {
	data := []byte(result.Name + "\n" + result.Key)
	fmt.Println("data", string(data))
	err := os.WriteFile(configreader.ConfigLocation+".lastfm_session", data, 0644)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(configreader.ConfigLocation + ".lastfm_session")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

}
