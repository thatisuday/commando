package main

import (
	"github.com/thatisuday/commando"
)

func main() {
	registry := commando.NewCommandRegistry()
	registry.SetExecutableName("")
	registry.Parse(nil)
}
