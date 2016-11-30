package main

import (
	"errors"
	"fmt"
	"strings"

	"bufio"
	"github.com/fatih/color"
	"github.com/jason0x43/go-toggl"
	"github.com/jinzhu/now"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"time"
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
	if len(c.Args()) < 1 {
		return startContinuousTimer(c)
	}

	session, err := initSession()
	if err != nil {
		return err
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

func startContinuousTimer(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return err
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	start := time.Now()

	event_queue := make(chan termbox.Event)
	go func() {
		for {
			event_queue <- termbox.PollEvent()
		}
	}()

loop:
	for {
		select {
		case ev := <-event_queue:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			default:
			}
		default:
		}
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		d := time.Now().Sub(start)
		timer := now.BeginningOfDay().Add(d)

		printCenter(1, "Time Entry", termbox.ColorYellow|termbox.AttrUnderline, termbox.ColorDefault)
		printCenter(3, timer.Format("15:04:05"), termbox.ColorWhite|termbox.AttrBold, termbox.ColorDefault)
		printCenter(4, fmt.Sprintf("%s %.2f", cfg.Currency, d.Hours()*cfg.Rate), termbox.ColorWhite, termbox.ColorDefault)
		printCenter(6, "(Press ESC to stop)", termbox.ColorMagenta, termbox.ColorDefault)
		termbox.Flush()
		time.Sleep(time.Millisecond)
	}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Close()
	reader := bufio.NewReader(os.Stdin)

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

	fmt.Printf("Timer stopped. You worked: %s\n", red(time.Now().Sub(start)))
	fmt.Printf("From %s to %s.\n", yellow(start.Format(time.Stamp)), yellow(time.Now().Format(time.Stamp)))
	fmt.Println("\nLet's add some info for that time entry before we send it.")

	var projs []toggl.Project
	for {
		projs, err = session.GetProjects()
		if err == nil {
			break
		}
		fmt.Println("Failed to retrieve project list. Are you connected?")
		time.Sleep(1 * time.Second)
		fmt.Println("Retrying...")
	}

	for i, proj := range projs {
		fmt.Printf("%d. %s \n", i+1, proj.Name)
	}
	fmt.Printf("Please select project number: ")

	var pindex int64
	for pindex == 0 {
		proj, _ := reader.ReadString('\n')
		pindex, err = strconv.ParseInt(strings.TrimSuffix(proj, "\n"), 10, 64)
		if err == nil && pindex > 0 && int(pindex) <= len(projs) {
			break
		}
		fmt.Printf("Invalid project, please select valid project from the list.\n Please select a project number: ")
	}

	fmt.Printf("Work description: ")
	description, _ := reader.ReadString('\n')

	pid := projs[pindex-1].ID

	_, err = session.CreateTimeEntry(start, time.Now(), strings.TrimSuffix(description, "\n"), pid)
	if err != nil {
		return fmt.Errorf("Failed to send time entry: %v", err)
	}

	fmt.Println("Time entry saved.")

	return nil
}

func printCenter(y int, s string, fg, bg termbox.Attribute) {
	w, _ := termbox.Size()
	x := w/2 - len(s)/2
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		x += w
	}
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
