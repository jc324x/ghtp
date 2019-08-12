// Package ght creates testable conditions in transient GitHub repos.
package ght

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/jychri/brf"
	"github.com/jychri/tilde"
)

// Read the contents of ~/.config/hub and verify its contents.
// Quit if the config file doesn't have a user value or OAuth token.
// The OAuth value is never assigned to a variable.
func readHubConfig() string {

	path := tilde.Abs("~/.config/hub")

	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("unable to open config at %v - %v", path, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var user string // GitHub username in ~/.config/hub
	var token bool  // OAuth token present in ~/.config/hub

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

type cleanup func()

// Create a directory that will hold all the repos.
// Returns a cleanup function that removes the temp directory.
func createTempDir(path string) cleanup {

	if path == "" {
		log.Fatalf("path cannot be ''")
	}

	path = tilde.Abs(path)

	os.RemoveAll(path)

	if err := os.MkdirAll(path, 0777); err != nil {
		log.Fatalf("Unable to create %v", path)
	}

	return func() {
		os.RemoveAll(path)
	}
}

type model struct {
	name   string // ght-Ahead
	remote string // jychri/ght-Ahead
	path   string // /Users/jychri/ght-testspace/ght-Ahead
}

type models []*model

func createModels(username string, path string, names []string) (models models) {
	path = tilde.Abs(path)

	for _, name := range names {
		model := new(model)
		model.name = name
		model.remote = strings.Join([]string{username, name}, "/")
		model.path = strings.Join([]string{path, name}, "/")
		models = append(models, model)
	}
	return models
}

// make directory
func (m *model) mkdir() {
	os.RemoveAll(m.path)
	os.MkdirAll(m.path, 0766)
}

// git init
func (m *model) init() {
	cmd := exec.Command("git", "init")
	cmd.Dir = m.path
	cmd.Run()
}

// hub delete -y m.remote; hub create
func (m *model) hub() {
	cmd := exec.Command("hub", "delete", "-y", m.remote)
	cmd.Dir = m.path
	cmd.Run()
	cmd = exec.Command("hub", "create")
	cmd.Dir = m.path
	cmd.Run()
}

// create a new file markdown file with Lorem Ipsum
func (m *model) create(name string) {
	lorem := "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
	data := []byte(lorem)
	name = strings.Join([]string{name, ".md"}, "")
	file := path.Join(m.path, name)

	if err := ioutil.WriteFile(file, data, 0777); err != nil {
		log.Fatal(err)
	}
}

// git add *
func (m *model) add() {
	cmd := exec.Command("git", "add", "*")
	cmd.Dir = m.path
	cmd.Run()
}

// git commit -m $message
func (m *model) commit(message string) {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = m.path
	cmd.Run()
}

// git push -u origin master
func (m *model) push() {
	cmd := exec.Command("git", "push", "-u", "origin", "master")
	cmd.Dir = m.path
	cmd.Run()
}

// Create temporary repos
func createTempRepos(models models, cleanup cleanup) {
	var wg sync.WaitGroup
	for i := range models {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			m.mkdir()                  // create model's directory
			m.init()                   // initialize a new Git repo
			m.hub()                    // create a fresh GitHub repo using hub
			m.create("README")         // create README.md with Lorem Ipsum
			m.add()                    // git add *
			m.commit("Initial commit") // git commit -m "Initial commit"
			m.push()                   // git push -u origin master
		}(models[i])
	}
	cleanup()
}

func cloneTempRepos(models models) {

	var wg sync.WaitGroup
	for i := range models {
		wg.Add(1)
		go func(m *model) {
			defer wg.Done()
			// m.clone()     // clone from GitHub
			// m.behind()    // set *Behind* models behind origin master
			// m.set(tdir)   // switch to tmpgis directory
			// m.ahead()     // set *Ahead* models behind origin master
			// m.untracked() // make *Untracked* models untracked
			// m.dirty()     // make *Dirty* models dirty
		}(models[i])
	}

}

// Temp ...
func Temp(path string, names []string) {
	username := readHubConfig()
	cleanup := createTempDir(path)
	models := createModels(username, path, names)
	createTempRepos(models, cleanup)
	cleanup = createTempDir(path)
	// stageTempRepos(models) // use path for 'overflow' => extra behind
	// cleanupTempDir = createTempDir(path)
}
