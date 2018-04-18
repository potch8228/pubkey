package pubkey

import (
	"bytes"
	"regexp"
	"testing"
)

func TestPubKey_LoadSettings(t *testing.T) {
	p := new(PubKey)
	p.newClient()
	p.loadFile("./settings_test.yml")

	if len(p.users) == 0 {
		t.Error("Failed to parse file")
	}

	for _, u := range p.users {
		if len(u.Id) == 0 {
			t.Error("Id data not parsed")
		}
	}
}

func TestPubKey_LoadPubKeys(t *testing.T) {
	p := new(PubKey)
	p.newClient()
	p.loadFile("./settings_test.yml")

	p.FillKeys()

	for _, v := range p.users {
		if len(v.Keys) == 0 {
			t.Error("Failed to get public key")
		}
	}
}

func TestPubKey_OutputList(t *testing.T) {
	p := new(PubKey)
	p.newClient()
	p.loadFile("./settings_test.yml")

	p.FillKeys()

	buf := new(bytes.Buffer)
	p.OutputList(buf)

	pattern := regexp.MustCompile(`^ssh.*\s{1}.*\n$`)
	if !pattern.Match(buf.Bytes()) {
		t.Error("Output format is not valid")
	}
}
