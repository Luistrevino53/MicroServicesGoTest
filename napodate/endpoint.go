package napodate

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

//Endpoint are exposed
type Endpoints struct {
	GetEndpoint			endpoint.Endpoint
	StatusEndpoint		endpoint.Endpoint
	ValidateEndpoint	endpoint.Endpoint
}

//MakeGetEndPoint returns the response from our service "get"
func MakeGetEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(getRequest) // we really just need the request, we don't use any value from it
		d, err := srv.Get(ctx)
		if err != nil{
			return getResponse{d, err.Error()}, nil
		}
		return getResponse{d, ""}, nil
	}
}

//MakeStatusEndpoint returns the response from our service "status"
func MakeStatusEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error){
		_ = request.(statusRequest) //we really just need the request, we don't use any value from it
		s, err := srv.Status(ctx)
		if err != nil {
			return statusResponse{s}, err
		}
		return statusResponse{s}, nil
	}
}

//MakeValidateEndpoint returns the response from our service "validate"
func MakeValidateEndpoint(srv Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (interface{}, error){
		req := request.(validRequest)
		b, err := srv.Validate(ctx, req.Date)
		if err != nil {
			return validResponse{b, err.Error()}, nil
		}
		return validResponse{b, ""}, nil
	}
}

//Get endpoint mapping 
func (e Endpoints) Get(ctx context.Context) (string, error){
	req := getRequest{}
	res, err := e.GetEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	getResp := res.(getResponse)
	if getResp.Err != "" {
		return "", errors.New(getResp.Err)
	}
	return getResp.Date, nil
}

//Status endpoint mapping
func (e Endpoints) Status(ctx context.Context) (string, error) {
	req := statusRequest{}
	resp, err := e.StatusEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	statusResp := resp.(statusResponse)
	return statusResp.Status, nil
}


//Validate endpoint mapping
func (e Endpoints) Validate(ctx context.Context, date string) (bool, error){
	req := validRequest{Date: date}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	validateResp := resp.(validResponse)
	if validateResp.Err != "" {
		return false, errors.New(validateResp.Err)
	}
	return validateResp.Valid, nil
}