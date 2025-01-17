// SPDX-FileCopyrightText: 2023 froggie <legal@frogg.ie>
//
// SPDX-License-Identifier: OSL-3.0

package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/effectindex/tripreporter/api"
	"github.com/effectindex/tripreporter/crypto"
	"github.com/effectindex/tripreporter/db"
	"github.com/effectindex/tripreporter/models"
	"github.com/effectindex/tripreporter/types"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var (
	dev    = flag.Bool("dev", false, "Run in development mode, alongside `make dev-ui`.")
	docker = flag.Bool("docker", false, "Run in Docker mode.")
)

func main() {
	flag.Parse()

	// Setup Zap for logging
	var err error
	var logger *zap.Logger

	if *dev {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sLogger := logger.Sugar()
	ctx := types.Context{Logger: sLogger}

	// Load and validate .env
	if err := godotenv.Load(); err != nil {
		logger.Fatal("err loading .env file (copy the .env.example)", zap.Error(err))
	}

	// Only SRV_ADDR, REDIS_PASS and DOCKER_SRV_HOST can be empty, anything else is validated.
	if err := validateEnvKeys(
		"SRV_PORT", "DEV_PORT", "SITE_NAME", "WORDLIST", "ACCOUNT_CONFIG", "DOCS_URL", "CORS_LOGGING", "VUE_APP_PROD_URL", "VUE_APP_DEV_URL", "VUE_APP_FORMKIT_API_KEY", "DB_NAME", "DB_USER", "DB_PASS", "DB_HOST", "DB_PORT", "REDIS_HOST", "REDIS_HOST", "DOCKER_POSTGRES_HOST", "DOCKER_REDIS_HOST",
	); err != nil {
		logger.Fatal("missing .env variables (copy the .env.example)", zap.Error(err))
	}

	// Setup NodeID for uuid generation
	if randomID, err := crypto.GenerateRandomBytes(6); err != nil {
		logger.Fatal("failed to initialize NodeID", zap.Error(err))
	} else {
		randomID[5] |= 0x01 // Set least significant bit of first true
		uuid.SetNodeID(randomID)
		ctx.Logger.Infof("Initialized random NodeID: %s", hex.EncodeToString(randomID))
	}

	// Setup required connections for postgresql and redis
	sDB := db.SetupDB(*docker, ctx.Logger)
	rDB := db.SetupRedis(*docker, ctx.Logger)

	defer sDB.Close()
	defer func(rDB *redis.Client) {
		err := rDB.Close()
		if err != nil {
			logger.Fatal("Failed to close Redis", zap.Error(err))
		}
	}(rDB)

	// Set context database now that we have one
	ctx.Database = sDB
	ctx.Cache = rDB

	// Setup wordlist and account configs
	models.SetupWordlist(ctx)
	models.SetupAccountConfig(ctx)

	// Setup proxy to webpack hot-reload server (for dev-ui) and regular http server (serves everything), and context
	api.Setup(*dev, ctx.Logger)
	api.SetupContext(ctx)
	api.SetupJwt()

	// Setup database patches right before server, now that everything else is ready
	db.SetupPatches(ctx)

	// Setup http server
	s := &http.Server{
		Addr:        os.Getenv("SRV_ADDR") + ":" + os.Getenv("SRV_PORT"),
		Handler:     api.CorsWrapper(api.Handler(), ctx.Logger),
		IdleTimeout: time.Minute,
	}

	if *dev {
		ctx.Logger.Infof("Running on %s in development mode...", s.Addr)
	} else {
		ctx.Logger.Infof("Running on %s in production mode...", s.Addr)
	}

	if err := s.ListenAndServe(); err != nil {
		ctx.Logger.DPanicf("Error in ListenAndServe: %v", err)
	}
}

func validateEnvKeys(keys ...string) error {
	missing := make([]string, 0)
	for _, key := range keys {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		return errors.New("[" + strings.Join(missing, ", ") + "]")
	}
	return nil
}
