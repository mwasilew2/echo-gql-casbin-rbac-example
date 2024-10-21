package graph

//go:generate go run generate.go

import (
	"log/slog"

	"gorm.io/gorm"

	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/service/authorization"
	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/util"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	db                   *gorm.DB
	logger               *slog.Logger
	ulidManager          *util.UlidManager
	authorizationService authorization.Authorization
}

func NewResolver(db *gorm.DB, logger *slog.Logger, ulidManager *util.UlidManager, authorizationService authorization.Authorization) *Resolver {
	logger = logger.With("subcomponent", "graph/Resolver")
	return &Resolver{
		db:                   db,
		logger:               logger,
		ulidManager:          ulidManager,
		authorizationService: authorizationService,
	}
}
