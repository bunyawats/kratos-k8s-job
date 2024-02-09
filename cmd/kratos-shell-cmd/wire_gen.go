// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-k8s-job/internal/biz"
	"kratos-k8s-job/internal/conf"
	"kratos-k8s-job/internal/data"
	"kratos-k8s-job/internal/server"
	"kratos-k8s-job/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, logLogger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logLogger)
	if err != nil {
		return nil, nil, err
	}
	jobRepo := data.NewJobRepo(dataData, logLogger)
	jobUseCase := biz.NewJobUseCase(jobRepo, logLogger)
	jobService := service.NewJobService(jobUseCase)
	grpcServer := server.NewGRPCServer(confServer, jobService, logLogger)
	httpServer := server.NewHTTPServer(confServer, jobService, logLogger)
	app := newApp(logLogger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
