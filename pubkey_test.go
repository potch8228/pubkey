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

func TestPubKey_OutputList(t *testing.T) {
	p := new(PubKey)
	p.newClient()
	p.loadFile("./settings_test.yml")

	p.FillKeys()

	buf := new(bytes.Buffer)
	p.OutputList(buf)

	pattern := regexp.MustCompile(`^ssh.*\s{1}.*$`)
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	for _, line := range lines {
		// last character is always "\n"
		if regexp.MustCompile(`^\s*$`).Match(line) {
			continue
		}

		if !pattern.Match(line) {
			t.Error("Output format is not valid")
		}
	}
}
