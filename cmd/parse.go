package cmd

import (
	"errors"
	"fmt"

	"code.gitea.io/gitea/modules/log"
	"github.com/Unknwon/com"
	"github.com/caiwp/ingest/modules/setting"
	"github.com/urfave/cli"
    "github.com/caiwp/ingest/modules/flume"
)

var CmdParse = cli.Command{
	Name:   "parse",
	Usage:  "Parse flume log and send to flume",
	Action: runParse,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "category, c",
		},
	},
}

func runParse(ctx *cli.Context) error {
	category := ctx.String("category")
	if category == "" {
		return errors.New("missing category")
	}
	if !com.IsSliceContainsStr(setting.ModelCategories, category) {
		log.Warn("category %s not found in app.ini", category)
		return errors.New(fmt.Sprintf("category %s not found in app.ini", category))
	}
	err := flume.Run(category)
	if err != nil {
		log.Error(4, "Parse %s failed: %v", category, err)
	}
	return err
}
