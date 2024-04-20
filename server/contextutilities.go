/*
   Package server provides functionality for creating and managing HTTP servers, including middleware support.

   Author: Sabyasachi Roy
*/

package server

import (
	"encoding/json"
	"io"
	"net/http"
)

// JSON writes a JSON response with the provided data to the ResponseWriter.
// It sets the Content-Type header to application/json.
// If an error occurs during JSON marshalling, it writes an error response with status code 500.
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

// DecodeJSON reads the JSON data from the request body and decodes it into the provided interface.
// It returns an error if reading or unmarshalling the JSON data fails.
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
