package kafka

import (
    "github.com/Shopify/sarama"
    "github.com/caiwp/ingest/modules/setting"
    "code.gitea.io/gitea/modules/log"
)

func SendMassage(s string) error {
    p := setting.Producer
    m := setting.Massage
    m.Value = sarama.ByteEncoder(s)
    partition, offset, err := p.SendMessage(m)
    if err != nil {
        return err
    }
    log.Trace("send massage return partition=%d, offset=%d", partition, offset)
    return nil
}