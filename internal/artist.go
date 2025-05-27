package internal

import (
	"strings"
	"log"
)

func GetArtist(track_info map[string]string) string {
	artist := track_info["Artist"]
	album_artist := track_info["AlbumArtist"]
	log.SetFlags(0)

	//check if artist field starts with album artist, overwrite if it does
	if strings.HasPrefix(artist, album_artist) {
		artist = album_artist
	} else {
		return SeparateArtists(artist)
	}
	return artist
}

func SeparateArtists(artist string) string {
	// list of artist separators to check for. could get from config file.
	//replace with regex??
	separators := []string {
		","," x ",";","/","&","feat.","Featuring","featuring",
	}
	// slice containing attempts to find 1st artist name
	attempts := []string{}

	// loop over the separators checking if any show up in the artist string
	for _, separator := range separators {
		artist, _, found := strings.Cut(artist, separator)
		if found {
			attempts = append(attempts, artist)
			// log.Print(separator)
			// log.Print("artist", artist)
			// log.Print("features", features)
		}
		// log.Print("no separator found!!")
	}
	return AttemptEval(attempts)
}

// pick one of the various attempts to move forward with
func AttemptEval(attempts []string) string {
	log.Print(attempts)
	match := strings.TrimSpace(attempts[0])

	return match
}
