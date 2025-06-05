package metadata

// this file solves the Main Problem i created this program to solve: submitting only the first artist in the metadata field
// to last.fm when provided with metadata in the form of a string separated by any number of arbitrary separators.
// this project began literally on day 1 of me learning go

import (
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"log"
	"regexp"
	"strings"
)

type Result struct {
	artist string
}

func GetArtist(artist string) string {
	config := configreader.ReadConfig()

	// this code is dumb. make it nicer later

	// if user opted out of first-artist scrobble, leave metadata as is
	if config.SingleArtist == false {
		return artist
	} else {
		// here we check if they want the metadata sanity check
		if config.SanityCheck == true {
			// if yes, run it.

			// this messes with reading rockbox log files. not doing it anymore, shouldnt be a problem anyway
			// step 1: check if first artist == albumArtist. easy
			// if strings.HasPrefix(artist, albumArtist) && len(albumArtist) > 0 {
			// 	return strings.TrimSpace(albumArtist)
			// } else
			if config.ApiCheck == true {
				// here we do the opt-out api-based check as a second-to-last resort
				return CheckMetadata(artist)
			} else {
				// if all else fails, run back to regex
				return separateArtists(artist)
			}
		} else {
			// if not, fall back regex-based cutting and allow the user to provide
			// a custom regex string in the config file
			return separateArtists(artist)
		}
	}
}

// parse result for artist info
func CheckMetadata(artist string) string {
	// split artist up across each separator then loop through them, popping the end off each time until a valid artist is found.
	artists := splitArtists(artist)
	for i := range artists {
		name := strings.Join(artists[:len(artists)-i], "")
		if SendQuery(name) != "Not an artist" {
			return strings.TrimSpace(name)
		}
	}
	return "failed to find artist"
}

func splitArtists(input string) []string {
	// log.Print("splitting artists")
	// Define a case-insensitive regex pattern for separators
	re := regexp.MustCompile(`(?i)\s*(,|;|&|feat\.|ft\.|featuring|and|\/)\s*`)

	// Split parts and find separators
	parts := re.Split(input, -1)
	separators := re.FindAllString(input, -1)

	// Build combined result
	var result []string
	for i, part := range parts {
		trimmedPart := strings.TrimSpace(part)
		if trimmedPart != "" {
			result = append(result, trimmedPart)
		}

		if i < len(separators) {
			sep := fmt.Sprintf("%v ", strings.TrimSpace(separators[i])) // capture group for the separator itself
			if sep != "" {
				result = append(result, sep)
			}
		}
	}
	return result
}

// this file needs a complete rewrite from here down lol. needs to support custom regex fields and also actually use regex
func separateArtists(artist string) string {
	// list of artist separators to check for. could get from config file.
	//replace with regex??
	// log.Print(artist)
	separators := []string{
		"feat.", "Featuring", "featuring", " x ", ",", ";", "/", "&", "and",
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
		} else {
			log.Print("no separator found!!")
		}
	}
	return attemptEval(attempts)
}

// pick one of the various attempts to move forward with - this is a stub and also a placeholder
func attemptEval(attempts []string) string {
	// log.Print(attempts)
	match := strings.TrimSpace(attempts[0])

	return match
}
