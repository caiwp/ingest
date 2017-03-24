package cron

import (
    "github.com/robfig/cron"
    "code.gitea.io/gitea/modules/log"
    "github.com/caiwp/ingest/modules/flume"
    "github.com/caiwp/ingest/modules/impala"
)

var (
    cr *cron.Cron
    ch chan struct{}
    running map[string]bool
)

func Run() {
    cr = cron.New()
    ch = make(chan struct{})

    cr.Start()
    defer cr.Stop()

    running = map[string]bool{
        "login": false,
    }

    cr.AddFunc("@every 2m", func() {
        s := "login"
        if running[s] {
            log.Info("Login is already running.")
            return
        }
        begin(s)
        defer end(s)

        flume.Run("login")

        impala.Run("login")

        return
    })

    <-ch
}

func begin(s string) {
    running[s] = true
    log.Info("%s start running.", s)
}

func end(s string) {
    running[s] = false
    log.Info("%s running end.", s)
}