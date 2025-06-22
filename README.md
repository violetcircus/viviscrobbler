# viviscrobbler
## what?
a last FM scrobbler written in Go, for MPD and Rockbox.

## why?
two main reasons.
1. no scrobbler exists that is exactly perfect for my use-case:

    a. I'm a recovering Spotify user, and the majority of my last.fm history is in the format spotify provides: only the first artist in is scrobbled when multiple are present on a song, and the rest are discarded. As far as I could tell from a cursory look, no other scrobbler for MPD has this behaviour - this presented a problem, as I wanted to have backwards compatibility with my old listens.

    a. I recently got an ipod 5 and put rockbox on it, and wanted the listening history from that to be backwards compatible too. 
2. i wanted to learn go, and figured this would be fun.

## how?
1. The scrobbler connects to MPD, then watches it for changes in the player subsystem.
2. If a new track starts, request the song info and player status from MPD once a second, and get the timestamp for when the song started.
3. Once the scrobble threshold (default 50%) has been reached, write the necessary information to create a scrobble to logFile.tsv
4. logFile.tsv is used as a queue, ensuring that even if the connection is dropped the scrobbles will be sent in the correct order.
5. The background thread checks the file for changes constantly, looping over any entries to send them to Last.FM.
6. By default, the artist name reported by MPD is adjusted to only contain the first artist if any features are present, using Musicbrainz or regex to check. The user can adjust this.
7. The MusicBrainz check works as follows: The artist string received from MPD is split along both "separators" - e.g., ",", "and", "Feat." etc. - and a slice - array in languages other than Go - is created containing the parts of the name and the separators in order. For example, "Tyler, the Creator Feat. Frank Ocean" would become ["Tyler", ",", "The Creator", "Feat.", "Frank Ocean"]. 
8. The program then concatenates these together and checks it against MusicBrainz's artist search API. If no artist matches this string, it removes the one at the end and checks that again, e.g., the previous example's second iteration would be "Tyler, The Creator Feat.". It repeats this process until it finds an exact match on Musicbrainz (though it is case-insensitive, and will return the version stored on MusicBrainz). This is a more sure-fire way of determining the first artist in a metadata segment than simple regex, as regex can easily be stumped by the presence of these separator characters within artist names - like, for example, "Tyler, The Creator" or "Earth, Wind & Fire". It is also better than simply using the Album Artist metadata field, as that is not always present and can often be misleading in the case of albums with songs created by multiple people, e.g. Trash Island by Drain Gang, or the OF Tape.
9. The resulting separated artist name is stored in mapFile.tsv along with the original string, removing the need for extra musicbrainz queries and allowing the user to customise the artist adjustments by simply editing the file.
10. The completed scrobble is sent to Last.FM.

To scrobble from a scrobble log created by Rockbox, you simply run the program with a file path as the argument. It will handle the artist metadata provided by the file in exactly the same way as it does the metadata received from MPD.

Tab-separated values files are used because tabs are the most reliable separator to use with data that can contain literally any punctuation character.
## installation:
either use the compiled release files or clone the repo and build it yourself - or, alternatively, use go install.
### dependencies:
[BurntSushi's toml parsing package](https://github.com/BurntSushi/toml)
### building from source:
```
to build from source run the following commands in your terminal:
git clone https://github.com/violetcircus/viviscrobbler
cd viviscrobbler
go build
```
this will create the viviscrobbler executable file, which you should move to your PATH.
### using go install
you can also install this with `go install github.com/violetcircus/viviscrobbler`.
## usage:
- `viviscrobbler`on its own will initiate the scrobbler.
- `viviscrobbler {PATH}` will load the file at the path and attempt to parse it for scrobbles - this is intended to be used with Rockbox's lastfm plugin's specific output format, so anything that doesn't conform to that format will not work!
- `viviscrobbler config` will regenerate the config file.
### files in the config folder:
- `config.toml` config file
- `logFile.tsv` tab-separated values file used as a queue for scrobbling. this ensures listens will be scrobbled in appropriate order even when offline
- `mapFile.tsv` tab-separated values file that stores the result of the program's artist-trimming along with the original string. You can edit this if you want to directly tell the program what to replace a certain artist metadata field with when scrobbling.
- `.lastfm_session` stores your lastfm session key generated after authorising the app along with your username for authentication purposes.
### systemd
downloading this application through any of these methods besides go install will provide a systemd service you can use which is pre-configured to only run after mpd starts.
### setup
IMPORTANT NOTE: If you build from source, you will need to provide your own API key and secret in the config file - which can be acquired from last.fm [here](https://www.last.fm/api/account/create). the release binaries, however, contain my api key and secret. 
#### Config options:
- serveraddress: string. address of the mpd server to connect to, defaults to localhost.
- serverport: string, port of the mpd server to connect to. defaults to 6600
- singleartist: bool. if true, will attempt to parse the artist section of a song's metadata for the first artist and provide only that to last.fm
- apicheck: bool. if true, will use musicbrainz's api to find the artist name rather than pure regex. if false, will use pure regex
- regex: string. if not blank, will be used as the regex string for separating artists (when api check is false)
- scrobblethreshold: what percentage of a song's duration needs to be listened to before it'll be sent to last.fm as a scrobble
- apikey: user API key
- secret: user API secret
### contributing:
- Not sure how active my maintaining of this will be, but I'll try to keep up with PRs. No promises on issues, though.
- The comments may not be exhaustive and I haven't written a single test as of yet, so it might be a bit annoying to parse my source code. This is also my first ever go project, so it's likely that it's very flawed.
- if you have suggestions, please create an issue.
### other information:
- Thanks to [YAMS](https://github.com/Berulacks/yams/) for being the main inspiration behind this project - their code was a great help while I was making this, and I directly lifted their systemd service, so go give them a star if you like this.
- This currently only works with Last.FM, because that's what I use. If you want to use another service I'm sure it wouldn't be hard to fork this and rework some of the api queries to point elsewhere.
- the scrobbler gets the first artist listed in metadata by splitting the artist string up across several separators and creating a slice consisting of each section - including the separators - then iterates over that slice, concatenating it together, checking that against musicbrainz's database, then dropping the end off and doing it again until it either finds an artist or reaches the beginning of the string.
### known issues:
- ~~for some reason, it keeps cutting artist names down further than it should - this seems to especially be an issue with Tyler, the Creator, who's been a thorn in the side of this program's development from the outset. Someone please tell artists to stop putting delimiters in their stage names.~~ solved this by reducing amount of spam to the musicbrainz api during artist checking via mapFile.tsv
