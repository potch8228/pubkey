package pubkey

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

const SETTINGS_FILE = "settings.yml"

type yamlSetting struct {
	Users []yamlUser
	Teams []yamlTeam
}

type yamlUser struct {
	Id string
}

type yamlTeam struct {
	Id string
}

type user struct {
	Id   string
	Keys []*github.Key
}

type PubKey struct {
	users  map[string]*user
	client *github.Client
}

// Initialize with custom setting file location
func NewPubKey() *PubKey {
	p := new(PubKey)
	p.newClient()
	p.load()
	return p
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

	s := yamlSetting{}
	if yaml.Unmarshal(d, &s) != nil {
		log.Fatalln("Can't unmarshal settings file")
	}

	p.users = make(map[string]*user)
	for i, _ := range s.Users {
		u := &user{s.Users[i].Id, nil}
		p.users[u.Id] = u
	}
	for i, _ := range s.Teams {
		t := s.Teams[i]
		users := p.GetMembers(t.Id)
		for i := range users {
			u := &user{users[i], nil}
			p.users[users[i]] = u
		}
	}
}

func (p *PubKey) GetMembers(team string) []string {
	users, resp, err := p.client.Organizations.ListMembers(context.Background(), team, nil)
	if err != nil {
		log.Println(resp)
		log.Fatalln("Failed to fetch members: " + team)
	}

	result := make([]string, len(users))
	for i := range users {
		u := users[i]
		result[i] = *u.Login
	}
	return result
}

func (p *PubKey) load() {
	p.loadFile(SETTINGS_FILE)
}

// NewClient creates APIClient
func (p *PubKey) newClient() {
	token := os.Getenv("GITHUB_TOKEN")
	var client *github.Client
	if token == "" {
		client = github.NewClient(nil)
	} else {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		client = github.NewClient(tc)
	}
	p.client = client
}

// Fetch user's public key data from GitHub
func (p *PubKey) FillKeys() *PubKey {
	for u, _ := range p.users {
		if u == "" {
			continue
		}
		keys, resp, err := p.client.Users.ListKeys(context.Background(), u, nil)
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
