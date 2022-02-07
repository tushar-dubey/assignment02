package app

import (
	"assignment02/internal/config"
	"context"
	"github.com/segmentio/kafka-go"
	"sync"
)

var (
	appContext *applicationContext
	once       sync.Once
)

type applicationContext struct {
	ctx context.Context

	config config.Config

	kafkaReader *kafka.Reader

	kafkaWriter *kafka.Writer
}

func NewAppContext(ctx context.Context) *applicationContext {
	once.Do(func() {
		appContext = new(applicationContext)
		appContext.ctx = ctx
	})
	return appContext
}

func GetAppContext() *applicationContext {
	return appContext
}

func (appContext *applicationContext) SetConfig(config config.Config) {
	appContext.config = config
}

func (appContext *applicationContext) SetKafkaReader(reader *kafka.Reader) {
	appContext.kafkaReader = reader
}

func (appContext *applicationContext) SetKafkaWriter(writer *kafka.Writer) {
	appContext.kafkaWriter = writer
}

func (appContext *applicationContext) GetKafkaReader() *kafka.Reader {
	return appContext.kafkaReader
}

func (appContext *applicationContext) GetKafkaWriter() *kafka.Writer {
	return appContext.kafkaWriter
}
