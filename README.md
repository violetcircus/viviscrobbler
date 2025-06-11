# viviscrobbler

## why?
two main reasons.
1. no scrobbler exists that is exactly perfect for my use-case:
    a. I'm a recovering Spotify user, and the majority of my last.fm history is in the format spotify provides: only the first artist in is scrobbled when multiple are present on a song, and the rest are discarded. As far as I could tell from a cursory look, no other scrobbler for MPD has this behaviour - this presented a problem, as I wanted to have backwards compatibility with my old listens.
    b. I recently got an ipod 5 and put rockbox on it, and wanted the listening history from that to be backwards compatible too. 
2. i wanted to learn go, and figured this would be fun.
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
you can also install this with `go install github.com/violetcircus/viviscrobbler@latest`. 
## usage:
`viviscrobbler`on its own will initiate the scrobbler.
`viviscrobbler {PATH}` will load the file at the path and attempt to parse it for scrobbles - this is intended to be used with Rockbox's lastfm plugin's specific output format, so anything that doesn't conform to that format will not work!
`viviscrobbler config` will regenerate the config file.
### systemd
downloading this application through any of these methods besides go install will provide a systemd service you can use which is pre-configured to only run after mpd starts.
### setup
IMPORTANT NOTE: If you build from source, you will need to provide your own API key and secret in the config file - which can be acquired from last.fm [here](https://www.last.fm/api/account/create). the release binaries, however, contain my api key and secret. 
#### Config options:
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
