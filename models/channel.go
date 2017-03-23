package models

import (
    "github.com/caiwp/ingest/modules/setting"
    "fmt"
)

type Channel struct {
    ChannelId int32
    ProductId int32
    Name      string
}

func GetChannel(pid int32, name string) (*Channel, error) {
    var c = new(Channel)
    query := fmt.Sprintf("SELECT channel_id, product_id, name FROM channels WHERE product_id = %d AND name = '%s' LIMIT 1", pid, name)
    q := setting.Db.NewQuery(query)
    err := q.One(c)
    if err != nil {
        return nil, err
    }
    if c.ChannelId == 0 {
        return nil, fmt.Errorf("product_id %d channel %s not found", pid, name)
    }
    return c, nil
}
