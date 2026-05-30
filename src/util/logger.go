package util

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type loggerContextKey string

var loggerKey = loggerContextKey("loggerKey")
var traceIDKey = loggerContextKey("traceIDKey")
var clientIPKey = loggerContextKey("clientIPKey")
var clientIDKey = loggerContextKey("clientIDKey")

// CustomLogFormatter ...
type CustomLogFormatter struct {
	Formatter logrus.Formatter
}

// GetLogger ...
func (f *CustomLogFormatter) GetLogger(writer io.Writer, app string,
	logLevel string) *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(writer)

	level := logrus.InfoLevel
	switch logLevel {
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	case "error":
		level = logrus.ErrorLevel
	case "info":
		level = logrus.InfoLevel
	}

	hostname, _ := os.Hostname()
	logger.SetLevel(level)

	fields := map[string]interface{}{
		"application": app,
		"logsource":   "application",
		"instance-id": hostname,
	}
	entry := logrus.NewEntry(logger).WithFields(fields)
	logger.SetFormatter(f)
	return entry
}

// NewCustomLogFormatter ...
func NewCustomLogFormatter() *CustomLogFormatter {
	jsonFormatter := new(logrus.JSONFormatter)
	jsonFormatter.TimestampFormat = "2006-01-02 15:04:05"
	return &CustomLogFormatter{
		Formatter: jsonFormatter,
	}
}

// Format renders a log entry.
func (f *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := logrus.Fields{}
	for key, value := range entry.Data {
		fields[key] = value
	}
	entry.Data = fields
	return f.Formatter.Format(entry)
}

// ServerEventFormatter ...
type ServerEventFormatter struct {
	Formatter logrus.Formatter
}

// GetLogger ...
func (f *ServerEventFormatter) GetLogger(writer io.Writer, app string,
	logLevel string) *logrus.Entry {

	logger := logrus.New()
	logger.SetOutput(writer)

	level := logrus.InfoLevel
	switch logLevel {
	case "debug":
		level = logrus.DebugLevel
	case "trace":
		level = logrus.TraceLevel
	case "error":
		level = logrus.ErrorLevel
	case "info":
		level = logrus.InfoLevel
	}

	hostname, _ := os.Hostname()
	logger.SetLevel(level)
	fields := map[string]interface{}{
		"application": app,
		"logsource":   "application",
		"instance-id": hostname,
	}
	entry := logrus.NewEntry(logger).WithFields(fields)
	logger.SetFormatter(f)
	return entry
}

// NewServerEventFormatter ...
func NewServerEventFormatter() *ServerEventFormatter {
	jsonFormatter := new(logrus.JSONFormatter)
	jsonFormatter.TimestampFormat = "2006-01-02 15:04:05"
	return &ServerEventFormatter{
		Formatter: jsonFormatter,
	}
}

// Format renders a log entry.
func (f *ServerEventFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fields := logrus.Fields{}
	for key, value := range entry.Data {
		fields[key] = value
	}

	if entry.Message != "" {
		handshakeErrString := "http: TLS handshake error from "
		if strings.Contains(entry.Message, handshakeErrString) {
			fields["secevent"] = LogSecEventServerFailure(entry.Message)
			clientIP := getClientIPFromLog(entry.Message, handshakeErrString, ":")
			if clientIP != "" {
				fields[LogClientIPFieldKey()] = clientIP
				fields[LogClientIDFieldKey()] = ""
			}
		}
	}

	entry.Data = fields
	return f.Formatter.Format(entry)
}

// Logger logger
func Logger(ctx context.Context) *logrus.Entry {
	if entry, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return entry
	}
	return nil
}

// WithLogger ...
func WithLogger(ctx context.Context, entry *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, entry)
}

// GenerateTraceID ...
func GenerateTraceID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// WithTraceID ...
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, traceIDKey, id)
}

// TraceID ...
func TraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		traceID = ""
	}
	return traceID
}

// WithClientIP ...
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

// ClientIP ...
func ClientIP(ctx context.Context) string {
	clientIP, ok := ctx.Value(clientIPKey).(string)
	if !ok {
		clientIP = ""
	}
	return clientIP
}

// WithClientID ...
func WithClientID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, clientIDKey, id)
}

// ClientID ...
func ClientID(ctx context.Context) string {
	clientID, ok := ctx.Value(clientIDKey).(string)
	if !ok {
		clientID = ""
	}
	return clientID
}

// HeaderTraceIDFieldKey ...
func HeaderTraceIDFieldKey() string {
	return "X-Correlation-ID"
}

// LogTraceIDFieldKey ...
func LogTraceIDFieldKey() string {
	return "x-correlation-id"
}

// LogOrderNumberFieldKey ...
func LogOrderNumberFieldKey() string {
	return "order_no"
}

// LogClientIPFieldKey ...
func LogClientIPFieldKey() string {
	return "x-client-ip"
}

// HeaderClientIPFieldKey ...
func HeaderClientIPFieldKey() string {
	return "X-Client-IP"
}

// LogClientIDFieldKey ...
func LogClientIDFieldKey() string {
	return "x-client-id"
}

// HeaderClientIDFieldKey ...
func HeaderClientIDFieldKey() string {
	return "X-Client-ID"
}

// LogSecEventServerFailure ...
func LogSecEventServerFailure(errString string) string {
	if (strings.Contains(errString, "TLS")) ||
		(strings.Contains(errString, "tls")) ||
		(strings.Contains(errString, "x509")) {
		return "MutualAuthenticationFailure"
	}
	return "HTTPServerFailure"
}

// LogSecEventClientFailure ...
func LogSecEventClientFailure(errString string, httpStatusCode int) string {
	if (strings.Contains(errString, "TLS")) ||
		(strings.Contains(errString, "tls")) ||
		(strings.Contains(errString, "x509")) {
		return "MutualAuthenticationFailure"
	}

	switch httpStatusCode {
	case http.StatusForbidden, http.StatusUnauthorized:
		return "ServerDeniedAccess"
	default:
		return "HTTPClientFailure"
	}
}

func getClientIPFromLog(message, start, end string) string {
	s := strings.Index(message, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(message[s:], end)
	if e == -1 {
		return ""
	}
	return message[s : s+e]
}
