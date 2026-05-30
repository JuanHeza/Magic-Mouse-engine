package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"magic-mouse-engine/middleware"
	"magic-mouse-engine/service"

	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

type contextKey string

const (
	ctxModelUpdateDigiRcptResponseKey     contextKey = "modelUpdateDigiRcptResponseKey"
	ctxModelSelectDigiRcptSaltResponseKey contextKey = "modelSelectDigiRcptSaltResponseKey"
	ctxModelGetOrderImageResponseKey      contextKey = "modelGetOrderImageResponseKey"
)

// MakeHTTPHandler : Create Http Request Handlers
func MakeHTTPHandler(s service.Service, logger log.Logger, csmw *middleware.ClientSetupMiddleware,
	log *logrus.Entry) http.Handler {

	r := mux.NewRouter()
	r.Use(csmw.SetClientSetupMiddleware)

	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}
	sr := r.PathPrefix("/kon-integration-digitalreceipt/api/v1/").Subrouter()

	sr.Methods("GET").Path("/healthz").Handler(httptransport.NewServer(
		e.HealthzEndPoint,
		decodeHealthzRequest,
		encodeHealthzResponse,
		options...,
	))

	sr.Methods("GET").Path("/version").Handler(httptransport.NewServer(
		e.VersionEndPoint,
		decodeVersionRequest,
		encodeVersionResponse,
		options...,
	))

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		url := fmt.Sprintf("%s%s", req.Host, req.RequestURI)
		log.WithFields(logrus.Fields{"event": "HTTPServer404"}).
			Error("Failed 404 on url: ", url)
		http.Error(w, "UPS! Parece ser que hubo un error. 404 Not found", http.StatusNotFound)
	})

	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	errCode, errMsg := codeAndMsgFrom(err)

	w.WriteHeader(errCode)
	err1 := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": errMsg,
	})
	if err1 != nil {
		panic("JSON encoding of error failure")
	}

}

func codeAndMsgFrom(err error) (int, string) {

	switch err {
	case service.ErrNotFound:
		return http.StatusNotFound, "Resource not found"
	case service.ErrBadRequest:
		return http.StatusBadRequest, "Invalid request"
	case service.ErrInternal:
		return http.StatusInternalServerError, "Internal server error"
	case service.ErrServiceUnavailable:
		return http.StatusServiceUnavailable, "Service is unavailable"
	case service.ErrGatewayTimeout:
		return http.StatusGatewayTimeout, "Gateway timeout"
	case service.ErrForbidden:
		return http.StatusForbidden, "Forbidden"
	case service.ErrUnauthorized:
		return http.StatusUnauthorized, "Unauthorized"
	case service.ErrConflict:
		return http.StatusConflict, "Conflict"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
