/*
   Package server provides functionality for creating and managing HTTP servers, including middleware support.

   Author: Sabyasachi Roy
*/

package server

// GetN registers a handler function for the GET method with the specified URL pattern.
// This method is intended to be used when creating new handler functions that accept a ContextHandler.
// If HandlerNew flag is set to false, log a fatal error message and return.
// The handler function should accept a ContextHandler as input.
func (api *MyAPIServer) GetN(pattern string, myHandler func(ctx ContextHandler)) {
	// Check if the new handler functions should be used
	if !api.HandlerNew {
		// If not, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function Get as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the GET method with the ServeMux
	api.Serv.ServeMux.HandleFunc("GET "+pattern, api.handlerWrapper(myHandler))
}

// PostN registers a handler function for the POST method with the specified URL pattern.
// This method is intended to be used when creating new handler functions that accept a ContextHandler.
// If HandlerNew flag is set to false, log a fatal error message and return.
// The handler function should accept a ContextHandler as input.
func (api *MyAPIServer) PostN(pattern string, myHandler func(ctx ContextHandler)) {
	// Check if the new handler functions should be used
	if !api.HandlerNew {
		// If not, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function Post as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the POST method with the ServeMux
	api.Serv.ServeMux.HandleFunc("POST "+pattern, api.handlerWrapper(myHandler))
}

// PutN registers a handler function for the PUT method with the specified URL pattern.
// This method is intended to be used when creating new handler functions that accept a ContextHandler.
// If HandlerNew flag is set to false, log a fatal error message and return.
// The handler function should accept a ContextHandler as input.
func (api *MyAPIServer) PutN(pattern string, myHandler func(ctx ContextHandler)) {
	// Check if the new handler functions should be used
	if !api.HandlerNew {
		// If not, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function Put as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the PUT method with the ServeMux
	api.Serv.ServeMux.HandleFunc("PUT "+pattern, api.handlerWrapper(myHandler))
}

// DeleteN registers a handler function for the DELETE method with the specified URL pattern.
// This method is intended to be used when creating new handler functions that accept a ContextHandler.
// If HandlerNew flag is set to false, log a fatal error message and return.
// The handler function should accept a ContextHandler as input.
func (api *MyAPIServer) DeleteN(pattern string, myHandler func(ctx ContextHandler)) {
	// Check if the new handler functions should be used
	if !api.HandlerNew {
		// If not, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function Delete as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the DELETE method with the ServeMux
	api.Serv.ServeMux.HandleFunc("DELETE "+pattern, api.handlerWrapper(myHandler))
}
