package main

import (
	"os"
	"time"

	"code.gitea.io/gitea/modules/log"
	"github.com/caiwp/ingest/cmd"
	"github.com/caiwp/ingest/modules/setting"
	"github.com/urfave/cli"
)

var Version = "0.1.0+dev"

func main() {
	var t = time.Now()

	app := cli.NewApp()
	app.Name = "Ingest"
	app.Usage = "数据流程服务"
	app.Version = Version
	app.Compiled = t
	app.Commands = []cli.Command{
		cmd.CmdParse,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Before = func(*cli.Context) error {
		setting.NewContext()
		setting.NewServices()
		return nil
	}
	app.After = func(*cli.Context) error {
		log.Warn("End time duration: %.4fs", time.Since(t).Seconds())
		setting.CloseServices()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Failed to run app with %s: %v", os.Args, err)
	}
}
