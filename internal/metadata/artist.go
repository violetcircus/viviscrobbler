package metadata

// this file solves the Main Problem i created this program to solve: submitting only the first artist in the metadata field
// to last.fm when provided with metadata in the form of a string separated by any number of arbitrary separators.

import (
	"fmt"
	"github.com/violetcircus/viviscrobbler/internal/configreader"
	"regexp"
	"strings"
)

type Result struct {
	artist string
}

func GetArtist(artist string) string {
	config := configreader.ReadConfig()

	// if user opted out of first-artist scrobble, leave metadata as is
	if config.SingleArtist == false {
		return artist
	} else {
		// here we check if they want the metadata sanity check
		if config.ApiCheck == true {
			// here we do the opt-out api-based check as a second-to-last resort
			return CheckMetadata(artist)
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
		fmt.Println("name:", name)
		if SendQuery(name) != "Not an artist" {
			return strings.TrimSpace(name)
		}
	}
	return separateArtists(artist)
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

// separate artists based on regex
func separateArtists(artist string) string {
	userRegex := configreader.ReadConfig().Regex
	var re *regexp.Regexp
	if userRegex != "" {
		re = regexp.MustCompile(userRegex)
	} else {
		re = regexp.MustCompile(`(?i)\s*(,|;|&|feat\.|ft\.|featuring|and|\/)\s*`)
	}
	parts := re.Split(artist, -1)
	return parts[0]
}
