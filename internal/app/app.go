package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/core-go/log"
	"github.com/core-go/mongo"
	"github.com/core-go/search"
	query "github.com/core-go/search/mongo"
	sv "github.com/core-go/service"
	v "github.com/core-go/service/v10"
	"reflect"

	. "go-service/internal/usecase/user"
)

type ApplicationContext struct {
	Health *health.Handler
	User   UserHandler
}

func NewApp(ctx context.Context, conf Config) (*ApplicationContext, error) {
	db, err := mongo.Setup(ctx, conf.Mongo)
	if err != nil {
		return nil, err
	}
	logError := log.ErrorMsg
	status := sv.InitializeStatus(conf.Status)
	action := sv.InitializeAction(conf.Action)
	validator := v.NewValidator()

	userType := reflect.TypeOf(User{})
	userQuery := query.UseQuery(userType)
	userSearchBuilder := mongo.NewSearchBuilder(db, "users", userQuery, search.GetSort)
	userRepository := mongo.NewRepository(db, "users", userType)
	userService := NewUserService(userRepository)
	userHandler := NewUserHandler(userSearchBuilder.Search, userService, status, logError, validator.Validate, &action)

	mongoChecker := mongo.NewHealthChecker(db)
	healthHandler := health.NewHandler(mongoChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
