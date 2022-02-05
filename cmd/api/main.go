package main

import (
	"assignment02/internal/boot"
	"context"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(boot.NewContext(context.Background()))
	defer cancel()

	// init app dependencies
	env := boot.GetEnv()
	err := boot.InitAPI(ctx, env)
	if err != nil {
		log.Fatalf("Failed to init API: %v", err)
	}
}
