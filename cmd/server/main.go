package main

import (
	"context"
	"github.com/jaro825/api-server-demo/cache"
	"github.com/jaro825/api-server-demo/server"
	"github.com/jaro825/api-server-demo/store"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"time"
)

var (
	configPath = pflag.String("config", "config.yaml", "configuration file")
	debugMode  = pflag.Bool("debug", false, "debug mode")
)

func main() {
	pflag.Parse()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if *debugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	config, err := LoadConfig(*configPath)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to load configuration")
	}

	dbStore, err := store.NewPostgresStore(logger, &store.Config{
		User:     config.DBUser,
		Password: config.DBPassword,
		Name:     config.DBName,
		Host:     config.DBHost,
		Port:     config.DBPort,
	})
	if err != nil {
		logger.Panic().Err(err).Msg("could not connect to a database")
	}

	redisCache, err := cache.NewRedisCache(config.RedisAddr, config.CacheTTL)
	if err != nil {
		logger.Panic().Err(err).Msg("could not connect to a redis cache")
	}

	defer redisCache.Stop()

	apiServer, err := server.New(logger, dbStore, redisCache, &server.Config{
		Port: config.APIPort,
	})
	if err != nil {
		logger.Panic().Err(err).Msg("could not create a service")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := apiServer.Run(); err != nil {
			logger.Panic().Err(err).Msg("server run error")
		}
	}()

	logger.Info().Msg("api server created")

	<-ctx.Done()

	shutdownCtx, shutdownStopCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownStopCtx()

	logger.Warn().Msg("received termination signal")

	if err := apiServer.Stop(shutdownCtx); err != nil {
		logger.Panic().Err(err).Msg("cannot gracefully stop the server")
	}

	logger.Info().Msg("application shutdown")
}
