package main

import (
	"assignment02/internal/boot"
	"assignment02/internal/worker"
	"context"
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
	err = boot.CreateKafkaTopic()
	if err != nil {
		log.Fatalf("Cannot create kafka topic: %v", err)
	}
	err = boot.InitKafkaReaderConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to kafka: %v", err)
	}
	worker.RunWorker(ctx)
}
