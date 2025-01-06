package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-kit/kit/transport/awslambda"
	internalevents "github.com/jeancaardo/go-app-event-notifier/internal/internal_events"
	"github.com/jeancaardo/go-app-event-notifier/pkg/bootstrap"
	"gorm.io/gorm"
)

var db *gorm.DB
var handler *awslambda.Handler

func init() {
	db = bootstrap.ConnectLocal()
	logger := bootstrap.InitSentry()

	eventsRepo := internalevents.NewRepository(db, logger)
	eventsSrv := internalevents.NewService(eventsRepo)
	re := internalevents.MakeEndpoints(eventsSrv)
	handler = internalevents.NewGetAllHandler(re, logger)
}

func main() {
	lambda.Start(handler)
}
