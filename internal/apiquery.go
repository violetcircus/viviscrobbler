package internal

import (
	"io"
	"net/url"
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"strings"
)

type Artist struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

type ArtistResponse struct {
	Artists []Artist `json:"artists"`
	Count int `json:"count"`
	Offset int `json:"offset"`
}

func SendQuery(artist string) string {
	// assemble the query URL
	query := fmt.Sprintf(`artist="%v"&limit=1`, artist)
	params := url.Values{}
	params.Add("query", query)
	params.Add("fmt", "json")

	finalUrl := fmt.Sprintf("https://musicbrainz.org/ws/2/artist/?%s", params.Encode())

	// send get request to URL
	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Println("whoops. fucked up on the get")
		log.Fatal(err)
	}
	// convert response into string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("whoops. fucked up on the io.readall")
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
	if len(artists) > limit{
		artists = artists[:limit]
	}
	// make that into a json object
	prettyJSON, err := json.MarshalIndent(artists, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(string(prettyJSON))

	// search artists for the artist name
	found := false 
	target := artist
	newArtist := ""
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(target)) {
			found = true
			newArtist = artist.Name
			break
		}
	}
	if !found {
		return "Not an artist"
	}
	return newArtist
}
