package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/go-slark/slark/logger"
	"github.com/go-slark/slark/pkg"
	"github.com/pkg/errors"
	"time"
)

type KafkaProducer struct {
	sarama.SyncProducer
	sarama.AsyncProducer
}

type ProducerConf struct {
	Brokers       []string `mapstructure:"brokers"`
	Retry         int      `mapstructure:"retry"`
	Ack           int16    `mapstructure:"ack"`
	ReturnSuccess bool     `mapstructure:"return_success"`
	ReturnErrors  bool     `mapstructure:"return_errors"`
}

type ConsumerGroupConf struct {
	Brokers      []string      `mapstructure:"brokers"`
	GroupID      string        `mapstructure:"group_id"`
	Topics       []string      `mapstructure:"topics"`
	Initial      int64         `mapstructure:"initial"`
	ReturnErrors bool          `mapstructure:"return_errors"`
	AutoCommit   bool          `mapstructure:"auto_commit"`
	Interval     time.Duration `mapstructure:"interval"`
	WorkerNum    uint          `mapstructure:"worker_num"`
}

type KafkaConf struct {
	Producer      *ProducerConf      `mapstructure:"producer"`
	ConsumerGroup *ConsumerGroupConf `mapstructure:"consumer_group"`
}

func (kp *KafkaProducer) Close() {
	_ = kp.SyncProducer.Close()
	kp.AsyncClose()
}

func (kp *KafkaProducer) SyncSend(ctx context.Context, topic, key string, msg []byte) error {
	traceID, ok := ctx.Value(utils.TraceID).(string)
	if !ok {
		traceID = utils.BuildRequestID()
	}
	pm := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
		Key:   sarama.StringEncoder(key),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(utils.TraceID),
				Value: []byte(traceID),
			},
		},
	}

	_, _, err := kp.SyncProducer.SendMessage(pm)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (kp *KafkaProducer) AsyncSend(ctx context.Context, topic, key string, msg []byte) error {
	traceID, ok := ctx.Value(utils.TraceID).(string)
	if !ok {
		traceID = utils.BuildRequestID()
	}
	pm := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msg),
		Key:   sarama.StringEncoder(key),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(utils.TraceID),
				Value: []byte(traceID),
			},
		},
	}

	kp.AsyncProducer.Input() <- pm
	select {
	case <-kp.AsyncProducer.Successes():

	case err := <-kp.AsyncProducer.Errors():
		return errors.WithStack(err)
	}
	return nil
}

func InitKafkaProducer(conf *ProducerConf) *KafkaProducer {
	return &KafkaProducer{
		SyncProducer:  newSyncProducer(conf),
		AsyncProducer: newAsyncProducer(conf),
	}
}

func newSyncProducer(conf *ProducerConf) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.RequiredAcks(conf.Ack) // WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Retry.Max = conf.Retry
	config.Producer.Return.Successes = conf.ReturnSuccess // true
	config.Producer.Return.Errors = conf.ReturnErrors     // true
	if err := config.Validate(); err != nil {
		panic(err)
	}

	producer, err := sarama.NewSyncProducer(conf.Brokers, config)
	if err != nil {
		panic(err)
	}
	return producer
}

func newAsyncProducer(conf *ProducerConf) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.RequiredAcks(conf.Ack) // WaitForAll
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Retry.Max = conf.Retry
	config.Producer.Return.Successes = conf.ReturnSuccess // true
	config.Producer.Return.Errors = conf.ReturnErrors     // true
	if err := config.Validate(); err != nil {
		panic(err)
	}

	producer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		panic(err)
	}

	return producer
}

type KafkaConsumerGroup struct {
	sarama.ConsumerGroup
	sarama.ConsumerGroupHandler
	Topics []string
	context.Context
	context.CancelFunc
	logger.Logger
}

func InitKafkaConsumer(conf *ConsumerGroupConf) *KafkaConsumerGroup {
	k := &KafkaConsumerGroup{
		ConsumerGroup: newConsumerGroup(conf),
		Topics:        conf.Topics,
	}
	k.Context, k.CancelFunc = context.WithCancel(context.TODO())
	return k
}

func newConsumerGroup(conf *ConsumerGroupConf) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = conf.Initial              // sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = conf.AutoCommit // false
	config.Consumer.Offsets.AutoCommit.Interval = conf.Interval * time.Millisecond
	config.Consumer.Return.Errors = conf.ReturnErrors // true
	if err := config.Validate(); err != nil {
		panic(err)
	}

	consumerGroup, err := sarama.NewConsumerGroup(conf.Brokers, conf.GroupID, config)
	if err != nil {
		panic(err)
	}

	return consumerGroup
}

func (kc *KafkaConsumerGroup) Consume() {
	for {
		err := kc.ConsumerGroup.Consume(kc.Context, kc.Topics, kc.ConsumerGroupHandler)
		if err != nil {
			kc.Log(kc.Context, logger.WarnLevel, map[string]interface{}{"error": err}, "consumer group consume fail")
		}
		if kc.Context.Err() != nil {
			kc.Log(kc.Context, logger.ErrorLevel, map[string]interface{}{"error": kc.Context.Err()}, "consumer group exit")
			return
		}
		time.Sleep(time.Second)
	}
}

func (kc *KafkaConsumerGroup) Start() error {
	kc.Consume()
	return nil
}

func (kc *KafkaConsumerGroup) Stop(_ context.Context) error {
	kc.CancelFunc()
	return kc.Close()
}

var (
	kafkaProducer      *KafkaProducer
	kafkaConsumerGroup *KafkaConsumerGroup
)

func GetKafkaProducer() *KafkaProducer {
	return kafkaProducer
}

func GetKafkaConsumerGroup() *KafkaConsumerGroup {
	return kafkaConsumerGroup
}
