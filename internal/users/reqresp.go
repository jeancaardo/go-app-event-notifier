package users

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/request"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/response"
)

type (
	GetReq struct {
		ID string `json:"id"`
	}

	StoreReq struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	UpdateReq struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	GetAllReq struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
		Sort  string `json:"sort"`
	}

	DeleteReq struct {
		ID string `json:"id"`
	}
)

func decodeStoreHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var user StoreReq
	if err := json.Unmarshal([]byte(event.Body), &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user, nil
}

func decodeUpdateHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var user UpdateReq
	if err := json.Unmarshal([]byte(event.Body), &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	if err := request.DecodeMap(event.PathParameters, &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user, nil
}

func decodeGetHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var user GetReq
	if err := request.DecodeMap(event.PathParameters, &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user, nil
}

func decodeGetAllHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var user GetAllReq
	if err := request.DecodeMap(event.QueryStringParameters, &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user, nil
}

func decodeDeleteHandler(ctx context.Context, payload []byte) (interface{}, error) {
	var event events.APIGatewayProxyRequest
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	var user DeleteReq
	if err := request.DecodeMap(event.PathParameters, &user); err != nil {
		return nil, response.BadRequest(err.Error())
	}

	return user, nil
}
