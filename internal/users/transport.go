package users

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/awslambda"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/lambda"
)

func NewStoreHandler(endpoints Endpoints, sentry log.Logger) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoints.Store,
		decodeStoreHandler,
		lambda.EncodeResponse,
		lambda.HandlerErrorEncoder(sentry),
		awslambda.HandlerFinalizer(lambda.HandlerFinalizer(sentry)))
}

func NewUpdateHandler(endpoints Endpoints, sentry log.Logger) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoints.Update,
		decodeUpdateHandler,
		lambda.EncodeResponse,
		lambda.HandlerErrorEncoder(sentry),
		awslambda.HandlerFinalizer(lambda.HandlerFinalizer(sentry)))
}

func NewGetHandler(endpoints Endpoints, sentry log.Logger) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoints.Get,
		decodeGetHandler,
		lambda.EncodeResponse,
		lambda.HandlerErrorEncoder(sentry),
		awslambda.HandlerFinalizer(lambda.HandlerFinalizer(sentry)))
}

func NewGetAllHandler(endpoints Endpoints, sentry log.Logger) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoints.GetAll,
		decodeGetAllHandler,
		lambda.EncodeResponse,
		lambda.HandlerErrorEncoder(sentry),
		awslambda.HandlerFinalizer(lambda.HandlerFinalizer(sentry)))
}

func NewDeleteHandler(endpoints Endpoints, sentry log.Logger) *awslambda.Handler {
	return awslambda.NewHandler(
		endpoints.Delete,
		decodeDeleteHandler,
		lambda.EncodeResponse,
		lambda.HandlerErrorEncoder(sentry),
		awslambda.HandlerFinalizer(lambda.HandlerFinalizer(sentry)))
}
