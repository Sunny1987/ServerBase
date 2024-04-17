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
	serveMux       *http.ServeMux
	middlewareList []Middleware
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
}

func (api *MyAPIServer) GetMyHttpServer() {
	api.Serv.serveMux = http.NewServeMux()
}

func (api *MyAPIServer) Get(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {

	api.Serv.serveMux.HandleFunc("GET /"+pattern, myHandler)
}

func (api *MyAPIServer) Post(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.serveMux.HandleFunc("GET /"+pattern, myHandler)
}
func (api *MyAPIServer) Put(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	api.Serv.serveMux.HandleFunc("GET /"+pattern, myHandler)
}

func (api *MyAPIServer) Prefix(prefix string) *http.ServeMux {
	v1 := http.NewServeMux()
	v1.Handle("/api/v1/", http.StripPrefix("/api/v1", api.Serv.serveMux))
	return v1
}

func (api *MyAPIServer) Run() error {
	var err error
	//log section
	l := log.New(os.Stdout, "MyAPIServer ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	//mux declaration and handler registrations
	api.GetMyHttpServer()

	//get registered middleware
	middlewareChain := api.MiddlewareChain(api.Serv.middlewareList)

	//get final middleware
	var servM *http.ServeMux
	if api.Serv.PrefixServeMux != nil {
		servM = api.Serv.PrefixServeMux
	} else {
		servM = api.Serv.serveMux
	}

	//Define server
	prodServer := &http.Server{
		Addr:         api.Addr,
		Handler:      middlewareChain(servM),
		ReadTimeout:  api.ReadTimeout,
		WriteTimeout: api.WriteTimeout,
		IdleTimeout:  api.IdleTimeout,
		ErrorLog:     l,
	}

	//call to serve
	go func() {
		myFigure := figure.NewFigure(api.AppName, "", true)
		myFigure.Print()
		l.Printf("version: %v", api.AppVer)
		l.Printf("Author: %v", api.AppAuthor)
		l.Printf("Starting server at port %v", api.Addr)
		if err = prodServer.ListenAndServe(); err != nil {
			l.Printf("Error starting server %v", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	sig := <-sigChan

	l.Println("Stopping server as per user interrupt", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = prodServer.Shutdown(tc)
	if err != nil {
		l.Println(err)
		return err
	}
	return err
}
