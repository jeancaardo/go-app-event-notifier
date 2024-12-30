package lambda

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/digitalhouse-tech/go-lib-kit/response"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/awslambda"
	"gorm.io/gorm"
)

// EncodeResponse
func EncodeResponse(_ context.Context, resp interface{}) ([]byte, error) {
	var res response.Response
	switch resp.(type) {
	case response.Response:
		res = resp.(response.Response)
	default:
		res = response.InternalServerError("unknown response type")
	}
	return APIGatewayProxyResponse(res)
}

// HandlerErrorEncoder
func HandlerErrorEncoder(log log.Logger) awslambda.HandlerOption {
	return awslambda.HandlerErrorEncoder(
		awslambda.ErrorEncoder(errorEncoder(log)),
	)
}

// HandlerFinalizer -
func HandlerFinalizer(log log.Logger) func(context.Context, []byte, error) {
	return func(ctx context.Context, resp []byte, err error) {
		if err != nil {
			log.Log("err", err)
		}
	}
}

func errorEncoder(log log.Logger) func(context.Context, error) ([]byte, error) {
	return func(_ context.Context, err error) ([]byte, error) {
		res := buildResponse(err, log)
		return APIGatewayProxyResponse(res)
	}
}

// buildResponse builds an error response from an error.
func buildResponse(err error, log log.Logger) response.Response {
	switch err.(type) {
	case response.Response:
		return err.(response.Response)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.NotFound("")
	}
	log.Log("err", err)
	return response.InternalServerError("")
}

// APIGatewayProxyResponse
func APIGatewayProxyResponse(res response.Response) ([]byte, error) {
	bytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	awsResponse := events.APIGatewayProxyResponse{
		Body:       string(bytes),
		StatusCode: res.StatusCode(),
		Headers:    res.GetHeaders(),
	}
	return json.Marshal(awsResponse)
}
