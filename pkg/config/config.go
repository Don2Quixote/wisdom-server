package config

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/env"
	"github.com/pkg/errors"
)

// ParseEnv is a wrapper around package confita to parse env variables into struct
// Example:
// var cfg struct {
//     Value string `config:"value,required"`
// }
// err = config.ParseEnv(&cfg)
//
// If no valuefor field with "required" struct tag - error will be returned.
func ParseEnv(cfg interface{}) error {
	err := confita.NewLoader(env.NewBackend()).Load(context.Background(), cfg)
	if err != nil {
		return errors.Wrap(err, "can't load env vars")
	}

	return nil
}
