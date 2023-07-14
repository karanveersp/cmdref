package main

import (
	"fmt"
	"log"

	"github.com/karanveersp/cmdref"
)

func main() {
	fileOps := cmdref.NewCmdFileOps()
	fmt.Printf("Commands file: %s\n", fileOps.GetFilePath())

	finished := false

	cmdMap, err := cmdref.LoadCommands(&fileOps)
	if err != nil {
		log.Fatal(err)
	}

	for !finished {
		action, err := cmdref.GetSelectedAction()
		if err != nil {
			log.Fatal(err)
		}
		switch action {
		case cmdref.Exit:
			fmt.Println("Bye!")
			finished = true
		default:
			cmdMap, err = cmdref.ProcessAction(cmdMap, action, &fileOps)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
