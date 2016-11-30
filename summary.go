package main

import (
	"fmt"

	"github.com/jinzhu/now"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
	"time"
)

func summary() cli.Command {
	return cli.Command{
		Name:   "summary",
		Usage:  "retrieve your account summary",
		Action: runSummary,
	}
}

var queryDateFormat = "2006-01-02"

func runSummary(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return err
	}

	report, err := session.GetSummaryReport(cfg.Workspace, time.Now().Format(queryDateFormat), time.Now().Format(queryDateFormat))
	if err != nil {
		return fmt.Errorf("Failed to load account summary: %s", err)
	}
	today := time.Duration(report.TotalGrand) * time.Millisecond

	report, err = session.GetSummaryReport(cfg.Workspace, now.BeginningOfWeek().Format(queryDateFormat), now.EndOfWeek().Format(queryDateFormat))
	if err != nil {
		return fmt.Errorf("Failed to load account summary: %s", err)
	}
	week := time.Duration(report.TotalGrand) * time.Millisecond

	report, err = session.GetSummaryReport(cfg.Workspace, now.BeginningOfMonth().Format(queryDateFormat), now.EndOfMonth().Format(queryDateFormat))
	if err != nil {
		return fmt.Errorf("Failed to load account summary: %s", err)
	}
	month := time.Duration(report.TotalGrand) * time.Millisecond

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	fmt.Fprint(w, "PERIOD", "\t", "HOURS", "\t", "INVOICE", "\n")
	fmt.Fprint(w, fmt.Sprintf("%s\t%.2f\t%.2f %s\n", "Today", today.Hours(), today.Hours()*cfg.Rate, cfg.Currency))
	fmt.Fprint(w, fmt.Sprintf("%s\t%.2f\t%.2f %s\n", "Week", week.Hours(), week.Hours()*cfg.Rate, cfg.Currency))
	fmt.Fprint(w, fmt.Sprintf("%s\t%.2f\t%.2f %s\n", "Month", month.Hours(), month.Hours()*cfg.Rate, cfg.Currency))
	w.Flush()

	return nil
}
