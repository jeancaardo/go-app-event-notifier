package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/transport/awslambda"
	"github.com/jeancaardo/go-app-event-notifier/internal/users"
	"github.com/jeancaardo/go-app-event-notifier/pkg/bootstrap"
	"gorm.io/gorm"
)

var db *gorm.DB
var handler *awslambda.Handler

func init() {
	db = bootstrap.ConnectLocal()
	logger := bootstrap.InitSentry()

	userRepo := users.NewRepository(db, logger)
	userSrv := users.NewService(userRepo)
	re := users.MakeEndpoints(userSrv)
	handler = users.NewGetAllHandler(re, logger)
}

func main() {
	lambda.Start(handler)
}
