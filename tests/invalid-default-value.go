package main

import (
	"github.com/thatisuday/commando"
)

func main() {
	registry := commando.NewCommandRegistry()
	registry.SetExecutableName("reactor")
	registry.Register("create").AddFlag("dir,d", "directory", commando.String, 21)
	registry.Parse(nil)
}
