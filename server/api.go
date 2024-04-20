/*
   Package server provides functionality for creating and running HTTP servers.

   Author: Sabyasachi Roy
*/

package server

import (
	"context"
	"github.com/common-nighthawk/go-figure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	AppNameDefault    = "MyAPIServer"
	AppVersionDefault = "1.0.0"
	AppAuthorDefault  = "Ghost"
)

// MyServer represents the server configuration and middleware.
type MyServer struct {
	// ServeMux is the default ServeMux for handling HTTP requests.
	ServeMux *http.ServeMux

	// MiddlewareList contains middleware functions to be applied to standard handler functions.
	MiddlewareList []Middleware

	// MiddlewareListN contains middleware functions to be applied to ContextHandler functions.
	MiddlewareListN []MiddlewareConvertedN

	// PrefixServeMux is an optional ServeMux for handling requests with a specific prefix.
	PrefixServeMux *http.ServeMux
}

// MyAPIServer represents the configuration for the API server.
type MyAPIServer struct {
	// Addr is the address the server will listen on.
	Addr string

	// Dns is the domain name of the server.
	Dns string

	// AppName is the name of the application.
	AppName string

	// AppVer is the version of the application.
	AppVer string

	// AppAuthor is the author of the application.
	AppAuthor string

	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum duration the server is allowed to idle without activity.
	IdleTimeout time.Duration

	// Serv is the instance of the MyServer.
	Serv *MyServer

	// Logger is the logger instance for logging server events.
	Logger *log.Logger

	// HandlerNew determines whether the server uses the new handler functions or the old ones.
	HandlerNew bool
}

// OptionalParams represents optional parameters for configuring the API server.
type OptionalParams struct {
	// Addr is the address the server will listen on.
	Addr string

	// Dns is the domain name of the server.
	Dns string

	// AppName is the name of the application.
	AppName string

	// AppVer is the version of the application.
	AppVer string

	// AppAuthor is the author of the application.
	AppAuthor string

	// ReadTimeout is the maximum duration for reading the entire request.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum duration the server is allowed to idle without activity.
	IdleTimeout time.Duration

	// Logger is the logger instance for logging server events.
	Logger *log.Logger

	// NewHandler determines whether the server uses the new handler functions or the old ones.
	NewHandler bool
}

// NewMyAPIServer creates a new instance of MyAPIServer with the provided optional parameters.
func NewMyAPIServer(opts *OptionalParams) *MyAPIServer {
	// Create a new MyAPIServer instance
	api := &MyAPIServer{}

	// Set port based on the provided options
	SetPort(opts, api)

	// Set DNS based on the provided options
	SetDNS(opts, api)

	// Set application name based on the provided options
	SetAppName(opts, api)

	// Set application version based on the provided options
	SetAppVersion(opts, api)

	// Set application author based on the provided options
	SetAppAuthor(opts, api)

	// Set read timeout based on the provided options
	SetReadTimeOut(opts, api)

	// Set write timeout based on the provided options
	SetWriteTimeOut(opts, api)

	// Set idle timeout based on the provided options
	SetIdleTimeOut(opts, api)

	// Set logger based on the provided options
	SetLogger(opts, api)

	// Set new handler flag based on the provided options
	SetNewHandler(opts, api)

	// Create a new MyServer instance with a new ServeMux
	api.Serv = &MyServer{ServeMux: http.NewServeMux()}

	return api
}

func SetNewHandler(opts *OptionalParams, api *MyAPIServer) {
	if opts.NewHandler == false {
		api.HandlerNew = false
	} else {
		api.HandlerNew = true
	}
}

func SetLogger(opts *OptionalParams, api *MyAPIServer) {
	if opts.Logger == nil {
		api.Logger = log.New(os.Stdout, api.AppName, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	} else {
		api.Logger = opts.Logger
	}
}

func SetIdleTimeOut(opts *OptionalParams, api *MyAPIServer) {
	if opts.IdleTimeout == 0 {
		api.IdleTimeout = 50 * time.Second
	} else {
		api.IdleTimeout = opts.IdleTimeout
	}
}

func SetWriteTimeOut(opts *OptionalParams, api *MyAPIServer) {
	if opts.WriteTimeout == 0 {
		api.WriteTimeout = 50 * time.Second
	} else {
		api.WriteTimeout = opts.WriteTimeout
	}
}

func SetReadTimeOut(opts *OptionalParams, api *MyAPIServer) {
	if opts.ReadTimeout == 0 {
		api.ReadTimeout = 20 * time.Second
	} else {
		api.ReadTimeout = opts.ReadTimeout
	}
}

func SetAppAuthor(opts *OptionalParams, api *MyAPIServer) {
	if opts.AppAuthor == "" {
		api.AppAuthor = AppAuthorDefault
	} else {
		api.AppAuthor = opts.AppAuthor
	}
}

func SetAppVersion(opts *OptionalParams, api *MyAPIServer) {
	if opts.AppVer == "" {
		api.AppVer = AppVersionDefault

	} else {
		api.AppVer = opts.AppVer
	}
}

func SetAppName(opts *OptionalParams, api *MyAPIServer) {
	if opts.AppName == "" {
		api.AppName = AppNameDefault
	} else {
		api.AppName = opts.AppName
	}
}

func SetDNS(opts *OptionalParams, api *MyAPIServer) {
	if opts.Dns == "" {
	} else {
		api.Dns = opts.Dns
	}
}

func SetPort(opts *OptionalParams, api *MyAPIServer) {
	if opts.Addr == "" {
		api.Addr = ":8080"
	} else {
		api.Addr = opts.Addr
	}
}

// ContextHandler wraps the response writer, request, logger, and DNS information.
type ContextHandler struct {
	// Writer is an interface used to construct HTTP responses.
	Writer http.ResponseWriter

	// Request is the HTTP request received from the client.
	Request *http.Request

	// Logger is a logger instance for logging context-related events.
	Logger *log.Logger

	// DNS is the domain name server information.
	DNS string
}

// handlerWrapper is a helper method that wraps a ContextHandler-based handler function into a standard http.HandlerFunc.
func (api *MyAPIServer) handlerWrapper(handler func(ContextHandler)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a ContextHandler with the request and response writer
		ctx := ContextHandler{
			Writer:  w,
			Request: r,
			Logger:  api.Logger,
			DNS:     api.Dns,
		}
		// Call the handler function with the ContextHandler
		handler(ctx)
	}
}

func (api *MyAPIServer) AddPrefix(prefix string) {
	v1 := http.NewServeMux()
	prefix2 := prefix[:len(prefix)-1]
	v1.Handle(prefix, http.StripPrefix(prefix2, api.Serv.ServeMux))
	api.Serv.PrefixServeMux = v1
}

func (api *MyAPIServer) Run() error {
	var err error
	var servM http.Handler

	if api.HandlerNew {
		servM = api.NewServMConfigure(servM)
	} else {
		servM = api.OldServMConfigure(servM)
	}
	api.Logger.Println("servM configured")

	//Define server
	prodServer := api.ConfigureServer(servM)
	api.Logger.Println("prodServer configured")

	//call to serve
	if err = api.StartServer(err, prodServer); err != nil {
		return err
	}
	sig := api.ListenForInterrupt()

	api.Logger.Println("Stopping server as per user interrupt", sig)
	err = api.ShutDown(err, prodServer)
	return err
}

func (api *MyAPIServer) ShutDown(err error, prodServer *http.Server) error {
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = prodServer.Shutdown(tc)
	if err != nil {
		api.Logger.Println(err)
		return err
	}
	return err
}

func (api *MyAPIServer) ListenForInterrupt() os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan
	return sig
}

func (api *MyAPIServer) StartServer(err error, prodServer *http.Server) error {
	go func() {
		myFigure := figure.NewFigure(api.AppName, "", true)
		myFigure.Print()
		api.Logger.Printf("version: %v", api.AppVer)
		api.Logger.Printf("Author: %v", api.AppAuthor)
		api.Logger.Printf("Starting server at port %v", api.Addr)
		if err = prodServer.ListenAndServe(); err != nil {
			api.Logger.Printf("Error starting server %v", err)
			os.Exit(1)
		}
	}()
	return err
}

func (api *MyAPIServer) ConfigureServer(servM http.Handler) *http.Server {
	prodServer := &http.Server{
		Addr:         api.Addr,
		Handler:      servM,
		ReadTimeout:  api.ReadTimeout,
		WriteTimeout: api.WriteTimeout,
		IdleTimeout:  api.IdleTimeout,
		ErrorLog:     api.Logger,
	}
	return prodServer
}

func (api *MyAPIServer) OldServMConfigure(servM http.Handler) http.Handler {
	//get registered middleware
	middlewareChain := api.MiddlewareChain(api.Serv.MiddlewareList)

	//get final middleware

	if api.Serv.PrefixServeMux != nil && api.Serv.MiddlewareList != nil {
		servM = middlewareChain(api.Serv.PrefixServeMux)
	} else if api.Serv.ServeMux != nil && api.Serv.MiddlewareList != nil {
		servM = middlewareChain(api.Serv.ServeMux)
	} else {
		servM = api.Serv.ServeMux
	}
	return servM
}

func (api *MyAPIServer) NewServMConfigure(servM http.Handler) http.Handler {
	middlewareChainN := api.MiddlewareChainN(api.Serv.MiddlewareListN)
	if api.Serv.PrefixServeMux != nil && api.Serv.MiddlewareListN != nil {
		servM = middlewareChainN(api.Serv.PrefixServeMux)
	} else if api.Serv.ServeMux != nil && api.Serv.MiddlewareListN != nil {
		servM = middlewareChainN(api.Serv.ServeMux)
	} else {
		servM = api.Serv.ServeMux
	}
	return servM
}
