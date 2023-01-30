//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"worker/dal"
	"worker/operator"
	"worker/service"
)

func InitContainer() (*Container, error) {
	wire.Build(
		// container
		NewContainer,

		// dal
		dal.NewMongoClient,

		// operator
		operator.NewMongoOperator,
		operator.NewStorageOperator,

		// service
		service.NewBucketService,
		service.NewFileService,
		service.NewUserService,
		service.NewHeartbeatService,
		service.NewPrometheusService,
	)
	return &Container{}, nil
}
