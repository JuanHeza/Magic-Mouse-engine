package service

import (
	"context"
	"errors"
	"io"

	"magic-mouse-engine/app"
	"magic-mouse-engine/client"
	"magic-mouse-engine/client/httpclient"
	"magic-mouse-engine/data"
	"magic-mouse-engine/model"

	"github.com/sirupsen/logrus"
)

// Service Api endpoint
type Service interface {
	Healthz(ctx context.Context) (model.HealthzResponseBody, error)
	Version(ctx context.Context) (model.VersionResponseBody, error)
}

// PosDigitalReceiptService Service Struct
type PosDigitalReceiptService struct {
	Config *data.Config
	Data   *data.Database
}

/*NewPosDigitalReceiptService ...*/
func NewPosDigitalReceiptService(config *data.Config, db *data.Database, c *httpclient.HTTPClient) Service {
	app.Create("/pos-receipt1", c, config)
	return &PosDigitalReceiptService{Config: config, Data: db}

}

// Declare errors
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

func ClientErrToServiceErr(err error) error {

	switch err {
	case client.ErrNotFound:
		return ErrNotFound
	case client.ErrInternal:
		return ErrInternal
	case client.ErrBadRequest:
		return ErrBadRequest
	case client.ErrServiceUnavailable:
		return ErrServiceUnavailable
	case client.ErrGatewayTimeout:
		return ErrGatewayTimeout
	case client.ErrForbidden:
		return ErrForbidden
	case client.ErrUnauthorized:
		return ErrUnauthorized
	case client.ErrConflict:
		return ErrConflict
	default:
		return ErrInternal
	}
}

/*POSService interface */
type POSService interface {
	GetPrefix() string
}

// Transaction represents an api instance
type Transaction struct {
	TempPrintFile func() (io.WriteCloser, error)
	Log           *logrus.Entry
}

// NewTransaction creates and return transaction object
func NewTransaction(tmpPrint func() (io.WriteCloser, error)) *Transaction {
	return &Transaction{
		TempPrintFile: tmpPrint,
	}
}
