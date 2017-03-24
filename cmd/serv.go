package cmd

import (
    "github.com/urfave/cli"
    "github.com/caiwp/ingest/modules/cron"
)

var CmdServ = cli.Command{
    Name:   "serv",
    Usage:  "Service the tasks",
    Action: runServ,
}

func runServ(*cli.Context) error {
    cron.Run()
    return nil
}