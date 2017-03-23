package models

import (
    "github.com/caiwp/ingest/modules/setting"
    "fmt"
)

type Gameserver struct {
    GameserverId int32
    ProductId    int32
    PlatformId   int32
    No           int32
}

func GetGameserver(pid, plid, no int32) (*Gameserver, error) {
    var g = new(Gameserver)
    query := fmt.Sprintf("SELECT " +
        "gameserver_id, product_id, platform_id, no " +
        "FROM gameservers " +
        "WHERE product_id = %d AND platform_id = %d AND no = %d " +
        "LIMIT 1", pid, plid, no)
    q := setting.Db.NewQuery(query)
    err := q.One(g)
    if err != nil {
        return nil, err
    }
    if g.GameserverId == 0 {
        return nil, fmt.Errorf("product_id %d platform_id %d gameserver_no %d not found", pid, plid, no)
    }
    return g, nil
}
