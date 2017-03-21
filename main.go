package main

import (
    "github.com/caiwp/ingest/modules/setting"
    "time"
    "code.gitea.io/gitea/modules/log"
    "github.com/caiwp/ingest/modules/kafka"
)

func main() {
    GlobalInit()
    t0 := time.Now()
    defer func() {
        log.Warn("End time duration: %.4fs", time.Since(t0).Seconds())
        setting.CloseServices()
    }()

    s := "hello world"
    err := kafka.SendMassage(s)
    if err != nil {
        log.Error(4, "send massage failed: %v", err)
    }
}

func GlobalInit() {
    setting.NewContext()
    setting.NewServices()
}