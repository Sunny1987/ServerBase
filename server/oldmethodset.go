/*
   Package server provides functionality for creating and managing HTTP servers, including middleware support.

   Author: Sabyasachi Roy
*/

package server

import "net/http"

// Get registers a handler function for the GET method with the specified URL pattern.
// If HandlerNew flag is set to true, use the new handler functions (GetN, PostN, etc.) instead.
// If HandlerNew flag is set to false, use the standard handler functions (Get, Post, etc.).
func (api *MyAPIServer) Get(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	// Check if the new handler functions should be used
	if api.HandlerNew {
		// If so, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function GetN as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the GET method with the ServeMux
	api.Serv.ServeMux.HandleFunc("GET "+pattern, myHandler)
}

// Post registers a handler function for the POST method with the specified URL pattern.
// If HandlerNew flag is set to true, use the new handler functions (PostN, GetN, etc.) instead.
// If HandlerNew flag is set to false, use the standard handler functions (Post, Get, etc.).
func (api *MyAPIServer) Post(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	// Check if the new handler functions should be used
	if api.HandlerNew {
		// If so, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function PostN as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the POST method with the ServeMux
	api.Serv.ServeMux.HandleFunc("POST "+pattern, myHandler)
}

// Put registers a handler function for the PUT method with the specified URL pattern.
// If HandlerNew flag is set to true, use the new handler functions (PutN, GetN, etc.) instead.
// If HandlerNew flag is set to false, use the standard handler functions (Put, Get, etc.).
func (api *MyAPIServer) Put(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	// Check if the new handler functions should be used
	if api.HandlerNew {
		// If so, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function PutN as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the PUT method with the ServeMux
	api.Serv.ServeMux.HandleFunc("PUT "+pattern, myHandler)
}

// Delete registers a handler function for the DELETE method with the specified URL pattern.
// If HandlerNew flag is set to true, use the new handler functions (DeleteN, GetN, etc.) instead.
// If HandlerNew flag is set to false, use the standard handler functions (Delete, Get, etc.).
func (api *MyAPIServer) Delete(pattern string, myHandler func(http.ResponseWriter, *http.Request)) {
	// Check if the new handler functions should be used
	if api.HandlerNew {
		// If so, log a fatal error message and return
		api.Logger.Fatalf("You need to call this function DeleteN as HandlerNew flag is set to %v", api.HandlerNew)
		return
	}
	// Register the handler function for the DELETE method with the ServeMux
	api.Serv.ServeMux.HandleFunc("DELETE "+pattern, myHandler)
}
