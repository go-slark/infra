package kafka

import "context"

type KafkaClient struct {
	*KafkaProducer
	*KafkaConsumerGroup
}

func (k *KafkaClient) AsyncProduce(ctx context.Context, topic, key string, msg []byte) error {
	return k.AsyncSend(ctx, topic, key, msg)
}

func (k *KafkaClient) SyncProduce(ctx context.Context, topic, key string, msg []byte) error {
	return k.SyncSend(ctx, topic, key, msg)
}

func (k *KafkaClient) Consume() error {
	go k.KafkaConsumerGroup.Consume()
	return nil
}

type queue interface {
	Produce(ctx context.Context, topic, key string, msg []byte) error
	Consume() error
}

func (k *KafkaClient) Produce(ctx context.Context, topic, key string, msg []byte) error {
	return k.AsyncSend(ctx, topic, key, msg)
}

var _ queue = &KafkaClient{}
