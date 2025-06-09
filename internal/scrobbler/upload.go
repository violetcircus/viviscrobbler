package scrobbler

import (
	"encoding/csv"
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

// upload scrobbles provided by whatever source: used by
// the thread that reads from the log and also when processing
// a user-provided file from rockbox
func UploadScrobbles(scrobble LoggedScrobble) bool {
	apiKey := secret.GetSecrets().ApiKey
	sk := getSession()

	// prepare the parameters for signature signing
	parameters := map[string]string{
		"artist":    scrobble.Artist,
		"track":     scrobble.Title,
		"album":     scrobble.Album,
		"timestamp": scrobble.Timestamp,
		"api_key":   apiKey,
		"sk":        sk,
		"method":    "track.scrobble",
	}
	signature := setup.SignSignature(parameters)
	// add signature to params after creating it
	parameters["api_sig"] = signature

	baseUrl := "https://ws.audioscrobbler.com/2.0/"

	// assemble params into a suitable format for post request
	postBody := url.Values{}
	for a, b := range parameters {
		path, err := url.PathUnescape(b)
		if err != nil {
			log.Fatal(err)
		}
		postBody.Set(a, path)
	}

	// send post request to scrobble api
	resp, err := http.Post(baseUrl, "application/x-www-form-urlencoded", strings.NewReader(postBody.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("scrobble response: %s", body)

	// return true if the scrobble was accepted: this will let the read scrobble loop know
	// to delete the first line
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	} else {
		return false
	}
}

// update last.fm's Now Playing api with current track - independent of scrobbling, so is run from main loop
// whenever the song changes
func UpdateNowPlaying(trackInfo metadata.TrackInfo) {
	apiKey := secret.GetSecrets().ApiKey
	sk := getSession()

	// prepare the parameters for signature signing
	parameters := map[string]string{
		"api_key": apiKey,
		"artist":  metadata.CheckMetadata(trackInfo.Artist),
		"track":   trackInfo.Title,
		"album":   trackInfo.Album,
		"method":  "track.updateNowPlaying",
		"sk":      sk,
	}
	signature := setup.SignSignature(parameters)
	// add signature to params after creating it
	parameters["api_sig"] = signature

	baseUrl := "https://ws.audioscrobbler.com/2.0/"

	// assemble params into a suitable format for post request
	postBody := url.Values{}
	for a, b := range parameters {
		path, err := url.PathUnescape(b)
		if err != nil {
			log.Fatal(err)
		}
		postBody.Set(a, path)
	}

	// send post request to now playing api
	resp, err := http.Post(baseUrl, "application/x-www-form-urlencoded", strings.NewReader(postBody.Encode()))
	if err != nil {
		log.Fatal(err)
	}
	_ = resp
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("now playing response: %s", body)
}
