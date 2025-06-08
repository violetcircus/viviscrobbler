package scrobbler

import (
	"encoding/csv"
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"github.com/violetcircus/viviscrobbler/internal/metadata"
	"github.com/violetcircus/viviscrobbler/internal/secret"
	"github.com/violetcircus/viviscrobbler/internal/setup"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// scrobble to be uploaded to lastfm
type Scrobble struct {
	Artist    string
	Album     string
	Timestamp string
	Title     string
	ApiKey    string
	Secret    string
	SKey      string
}

// get session key from session file
func getSession() string {
	f, err := os.Open(configreader.ConfigLocation + ".lastfm_session")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data[0][1]
}

func UploadScrobbles() {

}

func UpdateNowPlaying(trackInfo metadata.TrackInfo) {
	apiKey := secret.GetSecrets().ApiKey
	sk := getSession()

	parameters := map[string]string{
		"api_key": apiKey,
		"artist":  metadata.CheckMetadata(trackInfo.Artist),
		"track":   trackInfo.Title,
		"album":   trackInfo.Album,
		"method":  "track.updateNowPlaying",
		"sk":      sk,
	}
	signature := setup.SignSignature(parameters)
	parameters["api_sig"] = signature

	baseUrl := "https://ws.audioscrobbler.com/2.0/"

	postBody := url.Values{}
	for a, b := range parameters {
		path, err := url.PathUnescape(b)
		if err != nil {
			log.Fatal(err)
		}
		postBody.Set(a, path)
	}

	resp, err := http.Post(baseUrl, "application/x-www-form-urlencoded", strings.NewReader(postBody.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("now playing response: %s", body)
}
