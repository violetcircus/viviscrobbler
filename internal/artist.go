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
		artist, features, found := strings.Cut(artist, separator)
		if !found {
			log.Print("no separator found!!")
		}

		attempts = append(attempts, artist)
		log.Print(separator)
		log.Print("artist", artist)
		log.Print("features", features)
	}
	return artist
}
