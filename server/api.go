package server

import (
	"context"
	"encoding/json"
	"github.com/common-nighthawk/go-figure"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type MyServer struct {
	ServeMux       *http.ServeMux
	MiddlewareList []Middleware
	PrefixServeMux *http.ServeMux
}

type MyAPIServer struct {
	Addr         string
	Dns          string
	AppName      string
	AppVer       string
	AppAuthor    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Serv         *MyServer
	Logger       *log.Logger
	HandlerNew   bool
}

type OptionalParams struct {
	Addr         string
	Dns          string
	AppName      string
	AppVer       string
	AppAuthor    string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	Logger       *log.Logger
}

func NewMyAPIServer(opts *OptionalParams) *MyAPIServer {
	api := &MyAPIServer{}

	if opts.Addr == "" {
		api.Addr = ":8080"
	} else {
		api.Addr = opts.Addr
	}

	if opts.Dns == "" {
	} else {
		api.Dns = opts.Dns
	}

	if opts.AppName == "" {
		api.AppName = ""
	} else {
		api.AppName = opts.AppName
	}

	if opts.AppVer == "" {

	} else {
		api.AppVer = opts.AppVer
	}

	if opts.AppAuthor == "" {

	} else {
		api.AppAuthor = opts.AppAuthor
	}

	if opts.ReadTimeout == 0 {
		api.ReadTimeout = 20 * time.Second
	} else {
		api.ReadTimeout = opts.ReadTimeout
	}

	if opts.WriteTimeout == 0 {
		api.WriteTimeout = 50 * time.Second
	} else {
		api.WriteTimeout = opts.WriteTimeout
	}

	if opts.IdleTimeout == 0 {
		api.IdleTimeout = 50 * time.Second
	} else {
		api.IdleTimeout = opts.IdleTimeout
	}

	if opts.Logger == nil {
		api.Logger = log.New(os.Stdout, api.AppName, log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	} else {
		api.Logger = opts.Logger
	}

	api.Serv = &MyServer{ServeMux: http.NewServeMux()}

	return api
}

type ContextHandler struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Logger  *log.Logger
	DNS     string
}

// Define a wrapper function that converts ContextHandler into http.HandlerFunc
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

// ******Old Handler Definitions**********//
func (api *MyAPIServer) Get(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.ServeMux.HandleFunc("GET "+pattern, myHandler)
}
func (api *MyAPIServer) Post(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.ServeMux.HandleFunc("POST "+pattern, myHandler)
}
func (api *MyAPIServer) Put(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.ServeMux.HandleFunc("PUT "+pattern, myHandler)
}

func (api *MyAPIServer) Delete(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.ServeMux.HandleFunc("DELETE "+pattern, myHandler)
}

// ******New Handler Definitions**********//
func (api *MyAPIServer) GetN(pattern string, myHandler func(ctx ContextHandler)) {
	api.Serv.ServeMux.HandleFunc("GET "+pattern, api.handlerWrapper(myHandler))
}

func (api *MyAPIServer) PostN(pattern string, myHandler func(ctx ContextHandler)) {
	api.Serv.ServeMux.HandleFunc("POST "+pattern, api.handlerWrapper(myHandler))
}

func (api *MyAPIServer) PutN(pattern string, myHandler func(ctx ContextHandler)) {
	api.Serv.ServeMux.HandleFunc("PUT "+pattern, api.handlerWrapper(myHandler))
}

func (api *MyAPIServer) DeleteN(pattern string, myHandler func(ctx ContextHandler)) {
	api.Serv.ServeMux.HandleFunc("DELETE "+pattern, api.handlerWrapper(myHandler))
}

func (ctx *ContextHandler) JSON(data interface{}) {
	// Set Content-Type header to application/json
	ctx.Writer.Header().Set("Content-Type", "application/json")
	// Marshal the data to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		// If an error occurs during JSON marshalling, write an error response
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.Writer.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
	// Write the JSON response
	ctx.Writer.Write(jsonBytes)
}

// DecodeJSON decodes JSON data from the request body into the provided interface.
func (ctx *ContextHandler) DecodeJSON(v interface{}) error {
	// Read the request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	defer ctx.Request.Body.Close()

	// Unmarshal the JSON data into the provided interface
	if err = json.Unmarshal(body, &v); err != nil {
		return err
	}
	return nil
}

func (api *MyAPIServer) AddPrefix(prefix string) {
	v1 := http.NewServeMux()
	prefix2 := prefix[:len(prefix)-1]
	v1.Handle(prefix, http.StripPrefix(prefix2, api.Serv.ServeMux))
	api.Serv.PrefixServeMux = v1
}

func (api *MyAPIServer) Run() error {
	var err error

	//get registered middleware
	middlewareChain := api.MiddlewareChain(api.Serv.MiddlewareList)

	//get final middleware
	var servM http.Handler
	if api.Serv.PrefixServeMux != nil && api.Serv.MiddlewareList != nil {
		servM = middlewareChain(api.Serv.PrefixServeMux)
	} else if api.Serv.ServeMux != nil && api.Serv.MiddlewareList != nil {
		servM = middlewareChain(api.Serv.ServeMux)
	} else {
		servM = api.Serv.ServeMux
	}
	api.Logger.Println("servM configured")

	//Define server
	prodServer := &http.Server{
		Addr:         api.Addr,
		Handler:      servM,
		ReadTimeout:  api.ReadTimeout,
		WriteTimeout: api.WriteTimeout,
		IdleTimeout:  api.IdleTimeout,
		ErrorLog:     api.Logger,
	}

	api.Logger.Println("prodServer configured")

	//call to serve
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan

	api.Logger.Println("Stopping server as per user interrupt", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = prodServer.Shutdown(tc)
	if err != nil {
		api.Logger.Println(err)
		return err
	}
	return err
}
