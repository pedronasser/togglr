package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"net/http"

	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/jinzhu/now"
	"github.com/urfave/cli"
	"strconv"
)

type Invoice struct {
	From    string        `json:"from"`
	To      string        `json:"to"`
	Number  int           `json:"number"`
	Date    string        `json:"date"`
	DueDate string        `json:"due_date"`
	Notes   string        `json:"notes"`
	Items   []InvoiceItem `json:"items"`
}

type InvoiceItem struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	UnitCost float64 `json:"unit_cost"`
}

func invoice() cli.Command {
	return cli.Command{
		Name:      "invoice",
		Usage:     "send an invoice",
		ArgsUsage: "project",
		Flags:     invoiceFlags(),
		Action:    runInvoice,
	}
}

func invoiceFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringSliceFlag{
			Name:  "email",
			Usage: "send to emails",
		},
		cli.StringSliceFlag{
			Name:  "projects",
			Usage: "projects ids",
		},
		cli.StringFlag{
			Name:  "num",
			Usage: "invoice number",
		},
		cli.BoolFlag{
			Name:  "detailed",
			Usage: "download detailed report",
		},
		cli.StringFlag{
			Name:  "client",
			Usage: "invoice to...",
		},
		cli.StringFlag{
			Name:  "since",
			Usage: "from date (eg. 2016-11-01)",
		},
		cli.StringFlag{
			Name:  "until",
			Usage: "to date (eg. 2016-11-30)",
		},
	}
}

func runInvoice(c *cli.Context) error {
	session, err := initSession()
	if err != nil {
		return fmt.Errorf("Failed to init session: %s", err)
	}

	pids := c.StringSlice("projects")
	projects := []int{}

	for _, pid := range pids {
		p, err := getProjectId(pid)
		if err != nil {
			continue
		}
		projects = append(projects, p)
	}

	var from, to time.Time

	if c.String("since") != "" {
		from, err = time.Parse(queryDateFormat, c.String("since"))
		if err != nil {
			return fmt.Errorf("Invalid time format: %s", from)
		}
	}
	if c.String("until") != "" {
		to, err = time.Parse(queryDateFormat, c.String("until"))
		if err != nil {
			return fmt.Errorf("Invalid time format: %s", to)
		}
	}

	if c.String("client") == "" {
		return errors.New("Missing --client flag")
	}
	client := c.String("client")

	var num int64
	if c.String("num") != "" {
		num, err = strconv.ParseInt(c.String("num"), 10, 64)
	}

	if cfg.Rate == 0 {
		return errors.New("You must configure your rate first: togglr config rate <number>")
	}

	if cfg.Name == "" {
		return errors.New(`You must configure your name first: togglr config name "<your name>"`)
	}

	var fromStr, toStr string
	if !from.IsZero() {
		fromStr = from.Format(queryDateFormat)
	} else {
		fromStr = now.BeginningOfMonth().Format(queryDateFormat)
	}

	if !to.IsZero() {
		toStr = to.Format(queryDateFormat)
	} else {
		toStr = now.EndOfMonth().Format(queryDateFormat)
	}

	report, err := session.GetSummaryReport(cfg.Workspace, fromStr, toStr)
	if err != nil {
		logrus.Fatal(err)
	}

	var amount time.Duration
	for _, project := range report.Data {
		if len(projects) > 0 {
			for _, proj := range projects {
				if proj == project.ID {
					amount += time.Duration(project.Time) * time.Millisecond
				}
			}
			continue
		}

		amount += time.Duration(project.Time) * time.Millisecond
	}

	fmt.Printf("Creating invoice for %.2f hours - Total: %.2f %s\n", amount.Hours(), amount.Hours()*cfg.Rate, cfg.Currency)

	if c.Bool("detailed") {
		fmt.Println("Downloading detailed report from Toggl...")
		data, err := session.DownloadDetailedReport(cfg.Workspace, fromStr, toStr)
		if err != nil {
			logrus.Fatal(err)
		}
		os.Remove(fmt.Sprintf("detailed-%d-%s-%s.pdf", cfg.Workspace, fromStr, toStr))
		ioutil.WriteFile(fmt.Sprintf("detailed-%d-%s-%s.pdf", cfg.Workspace, fromStr, toStr), data, os.ModePerm)
	}

	fmt.Println("Generating invoice...")
	form := Invoice{
		From:   cfg.Name,
		To:     client,
		Number: int(num),
		Notes:  fmt.Sprintf("Invoice for working %.1f hours for %s", amount.Hours(), client),
		Items: []InvoiceItem{
			{
				Name:     fmt.Sprintf("Working hours for %s", client),
				Quantity: amount.Hours(),
				UnitCost: cfg.Rate,
			},
		},
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(form)
	if err != nil {
		logrus.Fatal(err)
	}

	resp, err := http.Post("https://invoice-generator.com", "application/json; charset=utf-8", b)
	if err != nil {
		logrus.Fatal(err)
	}

	data2, _ := ioutil.ReadAll(resp.Body)
	os.Remove(fmt.Sprintf("invoice-%d-%s-%s.pdf", cfg.Workspace, fromStr, toStr))
	ioutil.WriteFile(fmt.Sprintf("invoice-%d-%s-%s.pdf", cfg.Workspace, fromStr, toStr), data2, os.ModePerm)
	fmt.Println("done")

	return nil
}
