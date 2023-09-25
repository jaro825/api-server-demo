package store

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func (config Config) Validate() error {
	err := validation.ValidateStruct(&config,
		validation.Field(&config.User, validation.Required, is.Alpha),
		validation.Field(&config.User, validation.Required, validation.Length(2, 100)),
		validation.Field(&config.Password, validation.Required, is.ASCII),
		validation.Field(&config.Password, validation.Required, validation.Length(2, 100)),
		validation.Field(&config.Name, validation.Required, is.Alphanumeric),
		validation.Field(&config.Name, validation.Required, validation.Length(2, 100)),
		validation.Field(&config.Host, validation.Required, is.Host),
		validation.Field(&config.Port, validation.Required, is.Port),
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
