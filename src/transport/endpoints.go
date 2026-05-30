package transport

import (
	"context"
	"magic-mouse-engine/middleware"
	"magic-mouse-engine/model"
	"magic-mouse-engine/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
)

// Endpoints Type
type Endpoints struct {
	HealthzEndPoint        endpoint.Endpoint
	VersionEndPoint        endpoint.Endpoint
}

// MakeServerEndpoints methods
func MakeServerEndpoints(s service.Service) Endpoints {
	return Endpoints{
		HealthzEndPoint:        makeHealthzEndPoint(s),
		VersionEndPoint:        makeVersionEndPoint(s),
	}
}

func makeVersionEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g, e := s.Version(ctx)
		return model.VersionResponse{Body: g, Err: e}, nil
	}
}

func makeHealthzEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		g, e := s.Healthz(ctx)
		return model.HealthzResponse{Body: g, Err: e}, nil
	}
}

func ErrorLogger(ctx context.Context, err error) {

	logger := logrus.New()
	l := logger.WithField("alert", "true")
	cID := ctx.Value(middleware.CtxCorrelationID)
	if cID != nil {
		l = l.WithField(middleware.CorrelationHeader, cID)
	}
	l.Errorln("Endpoint failed with error", err.Error())
}
