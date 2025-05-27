package internal

import (
	"strings"
	"log"
)

func GetArtist(artist string) string {
	log.SetFlags(0)
	separators := []string{
		",",";","/","&","feat.","Featuring","featuring",
	}
	attempts:= []string{}
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

func AttemptEval(attempts []string) string {
	match := attempts[0]

	return match
}
