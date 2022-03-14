package app

import (
	"context"
	"time"

	"wisdom/internal/wisdom"

	"wisdom/pkg/config"
	"wisdom/pkg/logger"

	"github.com/pkg/errors"
)

// Run runs app. If returned error is not nil, program exited
// unexpectedly and non-zero code should be returned (os.Exit(1) or log.Fatal(...)).
func Run(ctx context.Context, log logger.Logger) error {
	log.Info("staring app")

	var cfg appConfig

	err := config.ParseEnv(&cfg)
	if err != nil {
		return errors.Wrap(err, "can't parse env")
	}

	server := wisdom.NewServer(cfg.Port, wisdom.PoW{
		ComplexityFactor:   cfg.ComplexityFactor,
		MaxComplexity:      cfg.MaxComplexity,
		ComplexityDuration: time.Second * time.Duration(cfg.ComplexityDurationSeconds),
	}, log)

	err = server.Launch(ctx)
	if err != nil {
		return errors.Wrap(err, "server error")
	}

	log.Info("app finished")

	return nil
}
