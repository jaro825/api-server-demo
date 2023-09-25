package server

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Port string
}

func (config Config) Validate() error {
	err := validation.ValidateStruct(&config, validation.Field(&config.Port, validation.Required, is.Port))
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
