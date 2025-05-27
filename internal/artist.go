package internal
// this file solves the Main Problem i created this program to solve: submitting only the first artist in the metadata field
// to last.fm when provided with metadata in the form of a string separated by any number of arbitrary separators.
// this project began literally on day 1 of me learning go

import (
	"net/http"
	"net/url"
	"strings"
	"log"
)

type Result struct {
	artist string
}

func GetArtist(trackInfo map[string]string) string {
	artist := trackInfo["Artist"]
	// albumArtist := trackInfo["AlbumArtist"]
	log.SetFlags(0)
	config := ReadConfig()

	// if user opted out of first-artist scrobble, leave metadata as is
	if config.SingleArtist == false {
		return artist
	} else {
		// here we check if they want the online metadata sanity check
		if config.SanityCheck == true {
			// if yes, run it. 
			return CheckMetadata(trackInfo)
		} else {
			// if not, fall back regex-based cutting and allow the user to provide
			// a custom regex string in the config file
			return SeparateArtists(artist)
		}
	}
}

func CheckMetadata(trackInfo map[string]string) string {
	// conn, err := net.Dial("tcp", "https://musicbrainz.org/ws/2/")
	// if err != nil {
	// 	log.Fatal(err) // for now - in future this will just jump over to SeparateArtists()
	// }
	// defer conn.Close()
	// okay i need to figure out how im doing this now lol

	baseUrl := "https://musicbrainz.org/ws/2/"
	// construct appropriate URL using track info

	// send get request 

	// parse result for artist info
	artist := trackInfo["Artist"]
	return artist
}

func SeparateArtists(artist string) string {
	// list of artist separators to check for. could get from config file.
	//replace with regex??
	log.Print(artist)
	separators := []string{
		"feat.","Featuring","featuring"," x ",",",";","/","&",
	}
	// slice containing attempts to find 1st artist name
	attempts := []string{}

	// loop over the separators checking if any show up in the artist string
	for _, separator := range separators {
		artist, _, found := strings.Cut(artist, separator)
		if found {
			attempts = append(attempts, artist)
			log.Print(separator)
			// log.Print("artist", artist)
		} else {
			log.Print("no separator found!!")
		}
	}
	return AttemptEval(attempts)
}

// pick one of the various attempts to move forward with - this is a stub and also a placeholder
func AttemptEval(attempts []string) string {
	log.Print(attempts)
	match := strings.TrimSpace(attempts[0])

	return match
}
