package base

import (
    "github.com/oschwald/geoip2-golang"
    "net"
    "code.gitea.io/gitea/modules/log"
)

func GetRegion(s string, path string) (string, error) {
    db, err := geoip2.Open(path)
    if err != nil {
        return "", err
    }
    defer db.Close()

    ip := net.ParseIP(s)
    record, err := db.City(ip)
    if err != nil {
        return "", err
    }

    if len(record.Subdivisions) > 0 {
        r := record.Subdivisions[0].Names["zh-CN"]
        return r, nil
    }
    log.Debug("record subdivisions %v", record.Subdivisions)
    return "", nil
}
