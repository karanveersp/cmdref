package main

import (
	"cmdref"
	"fmt"
	"log"
)

func main() {
	cmdFile, err := cmdref.CmdFilePath()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cmdFile)
	return
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
