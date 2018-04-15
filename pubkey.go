package pubkey

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"gopkg.in/yaml.v2"
)

const SETTINGS_FILE = "settings.yml"

type user struct {
	Id   string
	Keys []*github.Key
}

type PubKey struct {
	transport *http.Transport
	users     map[string]*user
}

// Initialize with default settings.yml location
func NewPubKey() *PubKey {
	p := new(PubKey)
	p.load()
	p.setHttpTransport()
	return p
}

// Initialize with custom setting file location
func NewPubKeyWithSettings(name string) *PubKey {
	p := new(PubKey)
	p.loadFile(name)
	p.setHttpTransport()
	return p
}

func (p *PubKey) setHttpTransport() {
	p.transport = &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
}

func (p *PubKey) loadFile(filename string) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		log.Fatalln("Can't open settings file")
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln("Can't read settings file")
	}

	var u []user
	if yaml.Unmarshal(d, &u) != nil {
		log.Fatalln("Can't unmarshal settings file")
	}

	p.users = make(map[string]*user, len(u))
	for i, _ := range u {
		p.users[u[i].Id] = &u[i]
	}
}

func (p *PubKey) load() {
	p.loadFile(SETTINGS_FILE)
}

// Fetch user's public key data from GitHub
func (p *PubKey) FillKeys() *PubKey {
	cli := github.NewClient(nil)
	for u, _ := range p.users {
		keys, resp, err := cli.Users.ListKeys(context.Background(), u, nil)
		if err != nil {
			log.Fatalln("Failed to fetch public key: " + u)
		}
		p.users[u].Keys = keys
		resp.Body.Close()
	}
	return p
}

// Print fetched user list as authorized_keys format on standard out
func (p *PubKey) PrintList() int {
	return p.OutputList(os.Stdout)
}

// Print fetched user list as authorized_keys format on given writer
func (p *PubKey) OutputList(to io.Writer) int {
	sb := new(strings.Builder)
	for u, _ := range p.users {
		for i, _ := range p.users[u].Keys {
			k := p.users[u].Keys[i]
			sb.WriteString(fmt.Sprintf("%s %s\n", *k.Key, u))
		}
	}
	fmt.Fprint(to, sb.String())
	return 0
}
