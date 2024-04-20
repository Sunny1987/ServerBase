# ServerBase

ServerBase is a Go library that provides a boilerplate for building HTTP servers with ease. It includes features such as middleware support, route registration, and server management to streamline your server development process.

## Installation

To use ServerBase in your Go project, simply import it:

```go
go get github.com/Sunny1987/ServerBase
```

# Getting Started

## Creating a Server
To create a new API server instance, you can use the NewMyAPIServer function provided by ServerBase. Here's an example of how to create a server:

```go
app := server.NewMyAPIServer(&server.OptionalParams{
    Addr:      ":8080",
    AppName:   "MyApp",
    AppAuthor: "John Doe",
    AppVer:    "1.0.0",
})
```

## Registering Routes
You can register HTTP routes with various methods such as **Get**, **Post**, **Put**, and **Delete**. Here's an example of registering routes:

```go
app.Get("/ping", pingHandler)
app.Post("/api/resource", createResourceHandler)
```

## Adding Middleware
ServerBase supports middleware to intercept and preprocess HTTP requests. You can add middleware functions using the **AddMiddleware method**. Here's an example:

```go
app.AddMiddleware(authMiddleware)
app.AddMiddleware(loggingMiddleware)
```

## Running the Server
To start the server, simply call the **Run** method:

```go
err := app.Run()
if err != nil {
    log.Fatal(err)
}
```

# Example Usage
Here's an example of how to use **ServerBase** to create and run an API server:

```go
package main

import (
    "github.com/Sunny1987/ServerBase/server"
    "log"
)

func main() {
    // Create a new API server instance
    app := server.NewMyAPIServer(&server.OptionalParams{
        Addr:      ":8080",
        AppName:   "MyApp",
        AppAuthor: "John Doe",
        AppVer:    "1.0.0",
    })

    // Register routes
    app.Get("/ping", pingHandler)
    app.Post("/api/resource", createResourceHandler)

    // Add middleware
    app.AddMiddleware(authMiddleware)
    app.AddMiddleware(loggingMiddleware)

    // Run the server
    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
    // Handler logic for the ping route
}

func createResourceHandler(w http.ResponseWriter, r *http.Request) {
    // Handler logic for creating a resource
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware logic for authentication
        next.ServeHTTP(w, r)
    })
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Middleware logic for logging
        next.ServeHTTP(w, r)
    })
}
```
# Alternate Approach with New Handlers
Here's an alternate example of how to use **ServerBase** to create and run an API server keeping **NewHandler=true**:

```go

import (
    "github.com/Sunny1987/ServerBase/server"
    "log"
)

func main() {
    // Create a new API server instance
    app := server.NewMyAPIServer(&server.OptionalParams{
        Addr:      ":8080",
        AppName:   "MyApp",
        AppAuthor: "John Doe",
        AppVer:    "1.0.0",
        NewHandler: true,
    })

    // Register routes
    app.GetN("/ping", pingHandler)
    app.PostN("/resource", createResourceHandler)

    // Add middleware
    app.AddMiddlewareN(authMiddleware)
    app.AddMiddlewareN(loggingMiddleware)

    // Run the server
    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func pingHandler(ctx server.ContextHandlert) {
    // Handler logic for the ping route
}

func createResourceHandler(ctx server.ContextHandler) {
    // Handler logic for creating a resource
}

func authMiddleware(ctx server.ContextHandler) {
   // Middleware logic for authentication
}

func loggingMiddleware(ctx server.ContextHandler) {
   // Middleware logic for logging
}
```
