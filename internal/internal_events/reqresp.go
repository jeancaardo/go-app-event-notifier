package internalevents

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/request"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/response"
	"time"
)

type (
	GetReq struct {
		ID string `json:"id"`
	}

	StoreReq struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Category    string    `json:"category"`
		Date        time.Time `json:"date"`
		Location    string    `json:"location"`
	}

	UpdateReq struct {
		ID          string    `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Category    string    `json:"category"`
		Date        time.Time `json:"date"`
		Location    string    `json:"location"`
	}

	GetAllReq struct {
		Category string    `json:"category"`
		DateFrom time.Time `json:"date_from"`
		DateTo   time.Time `json:"date_to"`
		Page     int       `json:"page"`
		Limit    int       `json:"limit"`
		Sort     string    `json:"sort"`
	}

	DeleteReq struct {
		ID string `json:"id"`
	}
)

func decodeStoreHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var awsEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &awsEvent); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var event StoreReq
	if err := json.Unmarshal([]byte(awsEvent.Body), &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return event, nil
}

func decodeUpdateHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var awsEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &awsEvent); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var event UpdateReq
	if err := json.Unmarshal([]byte(awsEvent.Body), &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	if err := request.DecodeMap(awsEvent.PathParameters, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return event, nil
}

func decodeGetHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var awsEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &awsEvent); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var event GetReq
	if err := request.DecodeMap(awsEvent.PathParameters, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return event, nil
}

func decodeGetAllHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var awsEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &awsEvent); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var event GetAllReq
	if err := request.DecodeMap(awsEvent.QueryStringParameters, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return event, nil
}

func decodeDeleteHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var awsEvent events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &awsEvent); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var event DeleteReq
	if err := request.DecodeMap(awsEvent.PathParameters, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return event, nil
}
