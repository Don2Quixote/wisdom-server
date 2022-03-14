package main

import (
	"context"

	"wisdom/internal/app"

	"wisdom/pkg/graceful"
	"wisdom/pkg/logger"

	"github.com/pkg/errors"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log := logger.NewLogrus()

	graceful.OnShutdown(cancel)

	err := app.Run(ctx, log)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error running app"))
	}
}
