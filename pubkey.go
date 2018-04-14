package pubkey

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const URL = "https://github.com/"
const URL_SUFFIX = ".keys"
const SETTINGS_FILE = "settings.yml"

type user struct {
	Id   string
	Keys []string
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
		panic("Can't open settings file")
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		panic("Can't read settings file")
	}

	var u []user
	if yaml.Unmarshal(d, &u) != nil {
		panic("Can't unmarshal settings file")
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
	for u, _ := range p.users {
		cli := &http.Client{Transport: p.transport}
		resp, err := cli.Get(URL + u + URL_SUFFIX)
		if err != nil {
			log.Fatalln("Failed to fetch public key: " + u)
			continue
		}
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			b := new(strings.Builder)
			b.Write(body)
			keys := strings.Split(b.String(), "\n")
			p.users[u].Keys = keys
		}
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
			if len(k) == 0 {
				continue
			}
			sb.WriteString(fmt.Sprintf("%s %s\n", k, u))
		}
	}
	fmt.Fprint(to, sb.String())
	return 0
}
