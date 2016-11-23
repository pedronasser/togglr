package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "togglr"
	app.Usage = "a tool for toggl"
	app.Authors = []cli.Author{{Name: "Pedro Nasser"}}
	app.CommandNotFound = func(c *cli.Context, cmd string) { fmt.Fprintf(os.Stderr, "command not found: %v\n", cmd) }
	app.Commands = []cli.Command{
		login(),
		configCmd(),
		summary(),
		projects(),
		start(),
		stop(),
		invoice(),
	}
	app.Run(os.Args)
}
