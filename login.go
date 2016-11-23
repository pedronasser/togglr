package main

import (
	"fmt"
	"os"

	"bufio"
	"github.com/urfave/cli"
	"golang.org/x/crypto/ssh/terminal"
)

func login() cli.Command {
	return cli.Command{
		Name:   "login",
		Usage:  "login your account",
		Action: loginToggle,
	}
}

func loginToggle(c *cli.Context) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Your Toggl account: ")
	username, _ := reader.ReadString('\n')
	fmt.Printf("Password: ")
	bytePassword, _ := terminal.ReadPassword(0)
	password := string(bytePassword)

	_, err := doLogin(username, password)
	if err != nil {
		return err
	}

	fmt.Println("\nYou are successfully logged to your Toggl account.")

	return nil
}
