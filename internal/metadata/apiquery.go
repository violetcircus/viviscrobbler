package metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Artist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ArtistResponse struct {
	Artists []Artist `json:"artists"`
	Count   int      `json:"count"`
	Offset  int      `json:"offset"`
}

func SendQuery(artist string) string {
	fmt.Println("received artist:", artist)
	// assemble the query URL
	query := fmt.Sprintf(`artist:"%v"`, artist)
	params := url.Values{}
	params.Add("query", query)
	params.Add("fmt", "json")

	finalUrl := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/?%s", params.Encode())

	// send get request to URL
	// error handling needs to fall back to other methods of checking metadata later
	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Println("whoops. messed up on the get")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// convert response into string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("whoops. messed up on the io.readall")
		log.Fatal(err)
	}

	// convert the body string into ArtistResponse struct
	var result ArtistResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	// crop it down to 5 results
	limit := 5
	artists := result.Artists
	if len(artists) > limit {
		artists = artists[:limit]
	}
	//make that into a json object for debugging
	// prettyJSON, err := json.MarshalIndent(artists, "", " ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Print(string(prettyJSON))
	fmt.Println("musicbrainz status code:", resp.StatusCode)

	// search artists for the artist name
	var found bool
	target := artist
	newArtist := ""
	for _, name := range artists {
		if strings.Contains(strings.ToLower(name.Name), strings.ToLower(target)) {
			found = true
			newArtist = name.Name
			break
		}
	}
	if !found {
		return "Not an artist"
	}
	return newArtist
}
