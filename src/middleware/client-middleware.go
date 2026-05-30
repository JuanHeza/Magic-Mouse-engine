package middleware

import (
	"magic-mouse-engine/data"
	"magic-mouse-engine/util"
	"net/http"

	"github.com/sirupsen/logrus"
)

// ClientSetupMiddleware ...
type ClientSetupMiddleware struct {
	config *data.Config
	logger *logrus.Entry
}

// NewClientSetupMiddleware Creates ClientSetupMiddleware instance
func NewClientSetupMiddleware(config *data.Config, log *logrus.Entry) *ClientSetupMiddleware {
	return &ClientSetupMiddleware{config: config, logger: log}
}

// SetClientSetupMiddleware Setup request parameters for Client
func (csmv *ClientSetupMiddleware) SetClientSetupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log := csmv.logger

		/* Trace ID */
		loggerTraceID := request.Header.Get(util.HeaderTraceIDFieldKey())
		if loggerTraceID == "" {
			loggerTraceID = util.GenerateTraceID()
		}
		log = log.WithFields(logrus.Fields{util.LogTraceIDFieldKey(): loggerTraceID})

		/* Client ID */
		clientID := request.Header.Get(util.HeaderClientIDFieldKey())
		log = log.WithFields(logrus.Fields{util.LogClientIDFieldKey(): clientID})


		ctx := request.Context()
		ctx = util.WithTraceID(ctx, loggerTraceID)
		ctx = util.WithClientID(ctx, clientID)

		request = request.WithContext(ctx)

		next.ServeHTTP(writer, request)
	})
}
