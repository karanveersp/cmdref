package main

import (
	"cmdref"
	"cmdref/prompter"
	"fmt"
	"log"
	"os"
)

func main() {
	cmdFile, err := cmdref.CmdFilePath()
	fmt.Printf("Commands file: %s\n", cmdFile)
	if err != nil {
		log.Fatal(err)
	}
	finished := false

	cmdsProvider := func() ([]byte, error) {
		data, err := os.ReadFile(cmdFile)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	cmdMap, err := cmdref.ParseCommands(cmdsProvider)
	if err != nil {
		log.Fatal(err)
	}

	for !finished {
		if len(cmdMap) == 0 {
			fmt.Println("No commands stored")
			createCmd, err := prompter.PromptConfirm("Store new command?")
			if err != nil {
				log.Fatalf("prompt failed %v\n", err)
			}
			if createCmd {
				cmdMap, err = cmdref.CreateHandler(cmdMap)
				if err != nil {
					log.Fatalf("error while creating command - %v\n", err)
				}
				err = cmdref.UpdateFile(cmdFile, cmdMap)
				if err != nil {
					log.Fatalf("error while updating file - %v\n", err)
				}
			}
		} else {
			action, err := cmdref.GetSelectedAction()
			if err != nil {
				log.Fatal(err)
			}
			switch action {
			case cmdref.Exit:
				fmt.Println("Bye!")
				finished = true
			default:
				cmdMap, err = cmdref.ProcessAction(cmdFile, cmdMap, action)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	//validate := func(input string) error {
	//	_, err := strconv.ParseFloat(input, 64)
	//	if err != nil {
	//		return errors.New("invalid number")
	//	}
	//	return nil
	//}
	//
	//prompt := promptui.Prompt{
	//	Label:    "Number",
	//	Validate: validate,
	//}
	//
	//result, err := prompt.Run()
	//
	//if err != nil {
	//	fmt.Printf("Prompt failed %v\n", err)
	//	return
	//}
	//
	//fmt.Printf("You choose %q\n", result)
}
