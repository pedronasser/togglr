package main

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jason0x43/go-toggl"
	"github.com/urfave/cli"
	"strconv"
)

func projects() cli.Command {
	return cli.Command{
		Name:   "projects",
		Usage:  "list projects",
		Action: listProjects,
		Subcommands: cli.Commands{
			{
				Name:      "alias",
				Usage:     "create alias for a project",
				Action:    aliasProject,
				ArgsUsage: "project_id alias",
			},
		},
	}
}

func listProjects(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return fmt.Errorf("Failed to init session: %s", err)
	}

	projs, err := session.GetProjects()
	if err != nil {
		return errors.New("Failed to retrieve projects")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprint(w, "ID", "\t", "NAME", "\n")
	for _, proj := range projs {
		fmt.Fprint(w, proj.ID, "\t", proj.Name, "\n")
	}
	w.Flush()

	return nil
}

func getProjectId(projectID string) (int, error) {
	pid, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		for p, alias := range cfg.Aliases {
			if alias == projectID {
				pid = int64(p)
				break
			}
		}

		if pid == 0 {
			return 0, errors.New("Invalid project id")
		}
	}

	return int(pid), nil
}

func aliasProject(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return fmt.Errorf("Failed to init session: %s", err)
	}

	projs, err := session.GetProjects()
	if err != nil {
		return errors.New("Failed to retrieve projects")
	}

	if len(c.Args()) < 1 {
		return errors.New("Missing project id")
	}

	if len(c.Args()) < 2 {
		return errors.New("Missing alias")
	}

	pid, _ := strconv.ParseInt(c.Args()[0], 10, 64)

	var project *toggl.Project
	for _, proj := range projs {
		if proj.ID == int(pid) {
			project = &proj
			break
		}
	}

	if project == nil {
		return fmt.Errorf("Couldn't find project: %d", pid)
	}

	cfg.Aliases[project.ID] = c.Args()[1]
	saveConfig()

	fmt.Printf("Defined `%s` as alias for project %d\n", c.Args()[1], pid)

	return nil
}
