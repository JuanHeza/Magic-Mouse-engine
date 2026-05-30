package client

import (
	"context"
	"errors"
	"net/http"
	"time"

	"magic-mouse-engine/client/httpclient"
	"magic-mouse-engine/util"
)

// HTTPClient interface
type HTTPClient interface {
	Send(context.Context, []byte, string) (resp *http.Response, err error)
}

// ClientHelper Created for writing functions common across client services
type ClientHelper struct {
	Name   string
	Client httpclient.Client
}

// Send HTTP request to client service
func (clntSvc *ClientHelper) Send(ctx context.Context,
	body []byte,
	headers map[string]string,
	url,
	method string) (resp *http.Response, err error) {

	log := util.Logger(ctx)

	defer func(t1 time.Time) {
		log.Infof("%s apiResponseTime %s\n", clntSvc.Name, time.Since(t1))
	}(time.Now())

	log.Infoln("Executing ", method, url)
	if headers == nil {
		headers = make(map[string]string)
	}

	//headers[string(globals.CorrelationHeader)], _ = ctx.Value(globals.CorrelationHeader).(string)
	return clntSvc.Client.Send(headers, body, method, url)
}

var (
	ErrNotFound           = errors.New("resource Not found")
	ErrInternal           = errors.New("internal server error")
	ErrBadRequest         = errors.New("bad request")
	ErrServiceUnavailable = errors.New("service Unavailable")
	ErrGatewayTimeout     = errors.New("gateway Timeout")
	ErrForbidden          = errors.New("forbidden Resource")
	ErrUnauthorized       = errors.New("unauthorized Client")
	ErrConflict           = errors.New("conflict")
)
