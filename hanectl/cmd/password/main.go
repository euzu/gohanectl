package main

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"gohanectl/hanectl/auth"
	"os"
	"time"
)

func generatePassword(password string) {

	var hash string
	pwdchan := make(chan string)
	go func() {
		hash, _ = auth.HashPassword(password) // ignore error for the sake of simplicity
		pwdchan <- hash
	}()

	spinner := "|/-\\"
	spinnerIndex := 0

	fmt.Println("\r")
	for {
		select {
		case hashed := <-pwdchan:
			fmt.Print("\r")
			fmt.Println(hashed)
			return
		default:
			{
				fmt.Printf("\r \033[36m%c\033[m  ", spinner[spinnerIndex%len(spinner)])
				time.Sleep(150 * time.Millisecond)
				spinnerIndex += 1
			}
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) > 0 {
		generatePassword(argsWithoutProg[0])
	} else {
		fmt.Print("Enter password: ")
		if password, err := terminal.ReadPassword(0); err == nil {
			generatePassword(string(password))
		}
	}
}
