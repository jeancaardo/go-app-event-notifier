package internalevents

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/jeancaardo/go-app-event-notifier/pkg/domain"
	"github.com/jeancaardo/go-app-event-notifier/pkg/utils/response"
	"net/http"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	// Endpoints struct
	Endpoints struct {
		Get    endpoint.Endpoint
		GetAll endpoint.Endpoint
		Store  endpoint.Endpoint
		Update endpoint.Endpoint
		Delete endpoint.Endpoint
	}
)

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Store:  makeStoreEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		event, err := s.Get(ctx, req.ID)
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", event, nil, nil), nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetAllReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		event, err := s.GetAll(ctx, Filters{
			Category: req.Category,
			DateFrom: req.DateFrom,
			DateTo:   req.DateTo,
			Page:     req.Page,
			Limit:    req.Limit,
			Sort:     req.Sort,
		})
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", event, nil, nil), nil
	}
}

func makeStoreEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(StoreReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		event := domain.Event{
			Name:        req.Name,
			Description: req.Description,
			Category:    req.Category,
			Location:    req.Location,
			Date:        req.Date,
		}
		u, err := s.Store(ctx, event)
		if err != nil {
			return verifyError(err)
		}
		return response.Created("success", u, nil, nil), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(UpdateReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		event := domain.Event{
			ID:          req.ID,
			Name:        req.Name,
			Description: req.Description,
			Category:    req.Category,
			Location:    req.Location,
			Date:        req.Date,
		}
		u, err := s.Update(ctx, event)
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", u, nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		err := s.Delete(ctx, req.ID)
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", nil, nil, nil), nil
	}
}

func verifyError(err error) (interface{}, error) {
	switch {
	case errors.Is(err, ErrEventNotFound):
		return nil, response.NotFound(err.Error())
	case errors.Is(err, ErrEventNameAlreadyExists):
		return nil, response.BadRequest(err.Error())
	case errors.Is(err, ErrOnStoreEvent) || errors.Is(err, ErrOnUpdateEvent):
		return nil, response.InternalServerError(err.Error())
	default:
		return nil, response.InternalServerError(err.Error())
	}
}
