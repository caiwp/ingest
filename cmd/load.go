package cmd

import (
    "github.com/urfave/cli"
    "fmt"
    "github.com/Unknwon/com"
    "github.com/caiwp/ingest/modules/setting"
    "code.gitea.io/gitea/modules/log"
    "github.com/caiwp/ingest/modules/impala"
)

var CmdLoad = cli.Command{
    Name:   "load",
    Usage:  "load hdfs data to impala",
    Action: runLoad,
    Flags: []cli.Flag{
        cli.StringFlag{
            Name: "category, c",
        },
    },
}

func runLoad(ctx *cli.Context) error {
    category := ctx.String("category")
    if category == "" {
        return fmt.Errorf("missing category")
    }
    if !com.IsSliceContainsStr(setting.ModelCategories, category) {
        log.Warn("category %s not found in app.ini", category)
        return fmt.Errorf("category %s not found in app.ini", category)
    }
    err := impala.Run(category)
    if err != nil {
        log.Error(4, "Parse %s failed: %v", category, err)
    }
    return err
}