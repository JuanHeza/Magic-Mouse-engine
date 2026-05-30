package middleware

import (
	"context"
	"time"

	"magic-mouse-engine/model"
	"magic-mouse-engine/service"
	"magic-mouse-engine/util"

	"github.com/sirupsen/logrus"
)

type ContextKey struct{}

var (
	CorrelationHeader = "X-Correlation-ID"
	CtxCorrelationID  ContextKey
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(service.Service) service.Service

// LoggingMiddleware Set logger
func LoggingMiddleware(logger *logrus.Entry) Middleware {
	return func(next service.Service) service.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   service.Service
	logger *logrus.Entry
}

func (mw loggingMiddleware) Healthz(ctx context.Context) (model.HealthzResponseBody, error) {
	mw.logger = mw.logger.WithFields(logrus.Fields{"api": "/healthz",
		"method":                   "GET",
		util.LogTraceIDFieldKey():  util.TraceID(ctx),
		util.LogClientIPFieldKey(): util.ClientIP(ctx),
		util.LogClientIDFieldKey(): util.ClientID(ctx)})

	mw.logger.WithFields(logrus.Fields{"start": time.Now()}).Info("Started")

	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{"took": time.Since(begin).String()}).Info("Completed")
	}(time.Now())
	ctx = util.WithLogger(ctx, mw.logger)

	return mw.next.Healthz(ctx)
}

func (mw loggingMiddleware) Version(ctx context.Context) (model.VersionResponseBody, error) {
	mw.logger = mw.logger.WithFields(logrus.Fields{"api": "/version",
		"method":                   "GET",
		util.LogTraceIDFieldKey():  util.TraceID(ctx),
		util.LogClientIPFieldKey(): util.ClientIP(ctx),
		util.LogClientIDFieldKey(): util.ClientID(ctx)})

	mw.logger.WithFields(logrus.Fields{"start": time.Now()}).Info("Started")

	defer func(begin time.Time) {
		mw.logger.WithFields(logrus.Fields{"took": time.Since(begin).String()}).Info("Completed")
	}(time.Now())

	ctx = util.WithLogger(ctx, mw.logger)
	return mw.next.Version(ctx)
}
