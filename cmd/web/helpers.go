package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the client
func (app *application) serverError(w http.ResponseWriter, err error) {
	// debug.Stack() function returns a slice of bytes containing a stack trace for the current goroutine
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Sends a specific status code and corresponding description to the client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("The template %s does not exist", page)
		app.serverError(w, err)
	}

	buf := new(bytes.Buffer)

	// write to buffer to handle possible error
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// no error, write contents of buffer to http.ResponseWriter
	w.WriteHeader(status)
	buf.WriteTo(w)

}
