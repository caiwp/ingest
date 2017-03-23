package models

import (
    "github.com/caiwp/ingest/modules/setting"
    "fmt"
)

type Platform struct {
    PlatformId int32
    ProductId  int32
    Name       string
}

func GetPlatform(pid int32, name string) (*Platform, error) {
    var p = new(Platform)
    query := fmt.Sprintf("SELECT platform_id, product_id, name FROM platforms WHERE product_id = %d AND name = '%s' LIMIT 1", pid, name)
    q := setting.Db.NewQuery(query)
    err := q.One(p)
    if err != nil {
        return nil, err
    }
    if p.PlatformId == 0 {
        return nil, fmt.Errorf("product_id %d platform %s not found", pid, name)
    }
    return p, nil
}
