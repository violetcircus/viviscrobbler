package setup

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"io"
	"log"
	"net/http"
	"sort"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func Setup() {

	//todo
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
	log.Println("signature", result)
	return result
}

func GetToken() string {
	parameters := make(map[string]string)
	parameters["apiKey"] = secret.GetSecrets().ApiKey
	signature := SignSignature(parameters)
	baseUrl := "https://ws.audioscrobbler.com/2.0/"
	log.Println("api key", parameters["apiKey"])
	urlParams := fmt.Sprintf("?method=auth.gettoken&api_key=%v&%v&format=json", parameters["apiKey"], signature)

	log.Println(baseUrl + urlParams)
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
