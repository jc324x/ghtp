// Package hubt creates testable conditions in transient GitHub repos.
package main

// This doesn't mess with flags or JSON files. Used while testing
// other packages. Dump in some JSON as a string or []byte.

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jychri/brf"
	"github.com/jychri/tilde"
)

// 1 id the user in ~/.config/hub || quit
// 2 unmarshal JSON || quit
// 3 create base directory || quit
// 4 expand repos into models || quit

// read the contents of ~/.config/hub and verify its contents.
// quit if the config file doesn't have a user or oauth token.
func readHubConfig() string {

	path := tilde.Abs("~/.config/hub")

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("unable to open config at %v - %v", path, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var user string // GitHub username in ~/.config/hub
	var token bool  // OAuth token in ~/.config/hub

	for scanner.Scan() {
		l := scanner.Text()

		if match, _ := brf.After(l, "- user:"); match != "" {
			user = match
		}

		if match, _ := brf.After(l, "oauth_token:"); match != "" {
			token = true
		}
	}

	if user == "" {
		log.Fatalf("No user in ~/.config.hub")
	}

	if token == false {
		log.Fatalf("No ouath token value in ~/.config.hub")
	}

	return user
}

// unmarsmall unmarshalls []byte into Config
// If something goes wrong, stop the run.
func unmarshal(bs []byte) (c Config) {
	if err := json.Unmarshal(bs, &c); err != nil {
		log.Fatalf("Unable to unmarshal data\n")
	}

	return c
}

// Config holds unmrashalled JSON from a gisrc.json file.
type Config struct {
	Bundles []struct {
		Path  string `json:"path"`
		Zones []struct {
			User      string   `json:"user"`
			Remote    string   `json:"remote"`
			Workspace string   `json:"workspace"`
			Repos     []string `json:"repositories"`
		} `json:"zones"`
	} `json:"bundles"`
}

// Init returns unmarshalled data from gisrc.json.
func Init(f flags.Flags) (c Config) {
	bs := read(f)     // read the file at path f.Config
	c = unmarshal(bs) // unmarshal the data from file at path f.Config
	return c
}

func parseJSON() {

}

func main() {
	fmt.Println(readHubConfig())
}
