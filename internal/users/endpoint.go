package users

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
		user, err := s.Get(ctx, req.ID)
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", user, nil, nil), nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetAllReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		users, err := s.GetAll(ctx, Filters{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		})
		if err != nil {
			return verifyError(err)
		}
		return response.OK("success", users, nil, nil), nil
	}
}

func makeStoreEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(StoreReq)
		if !ok {
			return nil, response.BadRequest(http.StatusText(http.StatusBadRequest))
		}
		user := domain.User{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}
		u, err := s.Store(ctx, user)
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
		user := domain.User{
			ID:    req.ID,
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		}
		u, err := s.Update(ctx, user)
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
	case errors.Is(err, ErrUserNotFound):
		return nil, response.NotFound(err.Error())
	case errors.Is(err, ErrUserEmailAlreadyExists):
		return nil, response.BadRequest(err.Error())
	case errors.Is(err, ErrOnStoreUser) || errors.Is(err, ErrOnUpdateUser):
		return nil, response.InternalServerError(err.Error())
	default:
		return nil, response.InternalServerError(err.Error())
	}
}
