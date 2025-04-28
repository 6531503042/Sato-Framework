package cli

import (
	"fmt"
	"os"
)

func Run() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: sato g module <name>")
		return
	}

	command := os.Args[1]
	entity := os.Args[2]
	name := os.Args[3]

	if command != "g" {
		fmt.Println("Unknown command:", command)
		return
	}

	switch entity {
	case "module":
		GenerateModule(name)
	default:
		fmt.Println("Unknown entity:", entity)
	}
}
