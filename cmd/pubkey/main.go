package main

import (
	"os"

	"github.com/potch8228/pubkey"
)

func main() {
	os.Exit(pubkey.NewPubKey().FillKeys().PrintList())
}
