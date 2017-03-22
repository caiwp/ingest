package kafka

import (
	"code.gitea.io/gitea/modules/log"
	"github.com/Shopify/sarama"
	"github.com/caiwp/ingest/modules/setting"
)

func SendMassage(s string, topic string) error {
	p := setting.Producer
	m := setting.Massage
    m.Topic = topic
	m.Value = sarama.ByteEncoder(s)
	partition, offset, err := p.SendMessage(m)
	if err != nil {
		return err
	}
	log.Info("send massage return partition=%d, offset=%d", partition, offset)
	return nil
}
