package db

import (
	"github.com/effectindex/tripreporter/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func TestGet(uuid uuid.UUID, ctx models.Context) (*models.Account, error) {
	a := &models.Account{
		Context: ctx,
		Unique:  models.Unique{ID: uuid},
	}

	if err := a.Get(); err != nil {
		ctx.Logger.Warnw("Failed to get test account", "account", a, zap.Error(err))
		return a, err
	} else {
		ctx.Logger.Infow("Got test account", "account", a)
		return a, nil
	}
}

func TestDelete(uuid uuid.UUID, ctx models.Context) (*models.Account, error) {
	a := &models.Account{
		Context:  ctx,
		Unique:   models.Unique{ID: uuid},
		Password: "examplePword",
	}

	if err := a.Delete(); err != nil {
		ctx.Logger.Warnw("Failed to delete test account", "account", a, zap.Error(err))
		return a, err
	} else {
		ctx.Logger.Infow("Deleted test account")
		return a, nil
	}
}

func TestCreate(ctx models.Context) (*models.Account, error) {
	if username, err := models.Wordlist.Random(3); err != nil {
		return nil, err
	} else {
		a := &models.Account{
			Context:  ctx,
			Type:     "Account",
			Email:    "user@email.com",
			Username: username,
			Password: "examplePword",
		}

		if err := a.Post(); err != nil {
			ctx.Logger.Warnw("Failed to make test account", "account", a, zap.Error(err))
			return a, err
		}

		ctx.Logger.Infow("Created test account")
		return a, nil
	}
}