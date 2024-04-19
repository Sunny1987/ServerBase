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
	l            *log.Logger
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
	l            *log.Logger
	//Prefix       string
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

	if opts.l == nil {
		api.l = log.New(os.Stdout, "MyAPIServer ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	} else {
		api.l = opts.l
	}

	api.Serv = &MyServer{ServeMux: http.NewServeMux()}

	return api
}

func (api *MyAPIServer) Get(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.l.Println("Starting Get")
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
	api.l.Println("servM configured")

	//Define server
	prodServer := &http.Server{
		Addr:         api.Addr,
		Handler:      servM,
		ReadTimeout:  api.ReadTimeout,
		WriteTimeout: api.WriteTimeout,
		IdleTimeout:  api.IdleTimeout,
		ErrorLog:     api.l,
	}

	api.l.Println("prodServer configured")

	//call to serve
	go func() {
		myFigure := figure.NewFigure(api.AppName, "", true)
		myFigure.Print()
		api.l.Printf("version: %v", api.AppVer)
		api.l.Printf("Author: %v", api.AppAuthor)
		api.l.Printf("Starting server at port %v", api.Addr)
		if err = prodServer.ListenAndServe(); err != nil {
			api.l.Printf("Error starting server %v", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan

	api.l.Println("Stopping server as per user interrupt", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = prodServer.Shutdown(tc)
	if err != nil {
		api.l.Println(err)
		return err
	}
	return err
}
