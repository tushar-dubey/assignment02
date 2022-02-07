package main

import (
	"assignment02/internal/boot"
	"assignment02/internal/boot/app"
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(boot.NewContext(context.Background()))
	defer cancel()

	// init app dependencies
	env := boot.GetEnv()
	err := boot.InitConfigAndDB(ctx, env)
	if err != nil {
		log.Fatalf("Failed to init API: %v", err)
	}
	err = boot.InitKafkaWriterConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to kafka: %v", err)
	}
	defer func(KafkaWriter *kafka.Writer) {
		err := KafkaWriter.Close()
		if err != nil {
			log.Fatalf("Cannot Close Kafka Writer Connection: %v", err.Error())
		}
	}(app.GetAppContext().GetKafkaWriter())
	err = boot.WriteAPIKeysToDB(ctx)
	if err != nil {
		log.Fatalf("Failed to write API keys to DB: %v", err.Error())
	}
	boot.RegisterAndServe(ctx)
}
