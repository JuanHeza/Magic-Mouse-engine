package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"magic-mouse-engine/client/httpclient"
	"magic-mouse-engine/data"
	"magic-mouse-engine/middleware"
	"magic-mouse-engine/service"
	"magic-mouse-engine/transport"
	"magic-mouse-engine/util"

	stdlog "log"

	gklog "github.com/go-kit/log"
)

func main() {
	/* Parse command line */
	configFile := flag.String("conf", "", "Configuration File")
	flag.Parse()

	if len(*configFile) == 0 {
		println("Conf file name required\n")
		os.Exit(1)
	}

	/* Parse configuration file to config */
	config := data.NewConfig()
	if err := config.Init(*configFile); err != nil {
		println("Cannot initialize config from file:", *configFile)
		println("Error:", err)
		os.Exit(1)
	}
	if err := util.InitVersion(); err != nil { //Initialize version variables
		println("Error in initializing version")
	}
	log := util.NewCustomLogFormatter().GetLogger(os.Stderr, "kon-integration-digitalreceipt", config.Logging.Level)
	log.WithFields(logrus.Fields{"event": "pos_gotemplate_kon_start"}).Info("Application starts")

	/* Initialize database connection */
	db := data.NewDatabase(log)
	if err := db.Init(config); err != nil {
		log.Error("Cannot init db:", err.Error())
		os.Exit(1)
	}
	defer db.Close(config)

	client := &httpclient.HTTPClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	var svc service.Service
	{
		svc = service.NewPosDigitalReceiptService(config, db, client)
		svc = middleware.LoggingMiddleware(log)(svc)
	}

	/* Listen and serve HTTP requests */
	url, err := url.Parse(config.Server.ListenAddress)
	if err != nil {
		log.Error("Cannot parse listen address")
		return
	}

	host, port, _ := net.SplitHostPort(url.Host)
	if host == "0.0.0.0" {
		host = ""
	}

	/* Initialize client setup middleware */
	clientMiddleware := middleware.NewClientSetupMiddleware(config, log)

	listenAddress := fmt.Sprintf("%s:%s", host, port)
	var httpHandler http.Handler
	{
		httpHandler = transport.MakeHTTPHandler(svc, gklog.NewLogfmtLogger(os.Stderr), clientMiddleware, log)
		//httpHandler = transport.MakeHTTPHandler(svc, gklog.NewLogfmtLogger(os.Stderr), clientMiddleware, log)
	}
	if r, ok := httpHandler.(*mux.Router); ok {
		fmt.Println("Registered routes:")
		r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			path, _ := route.GetPathTemplate()
			methods, _ := route.GetMethods()
			fmt.Printf("%v %v\n", methods, path)
			return nil
		})
	}
	/* Use custom logger for http/https transport errors */
	serverLog := util.NewServerEventFormatter().GetLogger(os.Stderr, "kon-integration-digitalreceipt", config.Logging.Level)
	logrusWriter := serverLog.Writer()
	defer logrusWriter.Close()

	server := &http.Server{
		Addr:    listenAddress,
		Handler: httpHandler,
		// TLSConfig: serverTLSConfig,
		ErrorLog: stdlog.New(logrusWriter, "", 0),
	}
	server.SetKeepAlivesEnabled(false)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		log.Info("HTTPS transport listening on port:", listenAddress)
		errs <- server.ListenAndServe()
	}()

	log.Info("exit", <-errs)
}
