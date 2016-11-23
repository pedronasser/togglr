package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func start() cli.Command {
	return cli.Command{
		Name:      "start",
		Usage:     "start a timer",
		ArgsUsage: `project_id ["description"]`,
		Action:    startTimer,
	}
}

func startTimer(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return err
	}

	if len(c.Args()) < 1 {
		return errors.New("Missing project id")
	}

	pid, err := getProjectId(c.Args()[0])
	if err != nil {
		return err
	}

	var description string
	if len(c.Args()) > 1 {
		description = strings.Join(c.Args()[1:], " ")
	}

	_, err = session.StartTimeEntryForProject(description, int(pid))
	if err != nil {
		return errors.New("Failed to start timer")
	}

	fmt.Println("Timer started. Have a good work!")
	return nil
}

func stop() cli.Command {
	return cli.Command{
		Name:   "stop",
		Usage:  "stop the current timer",
		Action: stopTimer,
	}
}

func stopTimer(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return err
	}

	time, err := session.GetCurrentTimeEntry()
	if err != nil {
		return errors.New("Failed to get current time entry")
	}

	_, err = session.StopTimeEntry(time)
	if err != nil {
		return errors.New("Failed to stop time entry")
	}

	fmt.Println("Timer stopped. Have a good rest.")
	return nil
}
