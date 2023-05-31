package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/lusis/htmxtools"
	"golang.org/x/exp/slog"
)

//go:embed templates/*
var templateFS embed.FS

var templates *template.Template

func main() {
	subbed, err := fs.Sub(templateFS, "templates")
	if err != nil {
		panic(err)
	}
	tmpls, err := template.New("").ParseFS(subbed, "*")
	if err != nil {
		panic(err)
	}
	templates = tmpls
	mux := http.NewServeMux()
	mux.HandleFunc("/alert", serversideAlert)
	mux.HandleFunc("/button-push", buttonPush)
	mux.HandleFunc("/", templateMiddleware)
	middleware := htmxtools.Middleware(requestLogger(mux))
	if err := http.ListenAndServe(":3000", middleware); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func templateMiddleware(w http.ResponseWriter, r *http.Request) {
	// let's trim off the leftmost slash to find our template name
	p := strings.TrimLeft(r.URL.Path, "/")
	if p == "" {
		// server index by default
		p = "index.html"
	}
	htmxRequest := htmxtools.RequestFromContext(r.Context())
	data := struct {
		HtmxRequest *htmxtools.HTMXRequest
	}{
		HtmxRequest: htmxRequest,
	}

	if t := templates.Lookup(p); t != nil {
		if err := t.Execute(w, data); err != nil {
			slog.Error("unable to render template", "error", err)
			return
		}
	}
}

func serversideAlert(w http.ResponseWriter, r *http.Request) {
	if htmxRequest := htmxtools.RequestFromContext(r.Context()); htmxRequest != nil {
		// we're also going to ensure that this alert link doesn't make it's way into the history
		w.Header().Set(htmxtools.ReplaceURLResponse.String(), htmxRequest.CurrentURL)
		w.Header().Set(htmxtools.TriggerResponse.String(), `{"showMessage":{"level" : "info", "message" : "this alert was trigged via the HX-Trigger header"}}`)
	}
}
func buttonPush(w http.ResponseWriter, r *http.Request) {
	htmxRequest := htmxtools.RequestFromContext(r.Context())
	if htmxRequest != nil {
		j, err := json.Marshal(htmxRequest)
		if err != nil {
			w.Write([]byte("unable to encode htmxrequest to json: " + err.Error())) // nolint: errcheck
		}
		response := `
		<div id="button-push-response">Your request:<pre>%s</pre></div>
		`
		w.Write([]byte(fmt.Sprintf(response, j))) // nolint: errcheck
	}
}

func requestLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		htmxreq := htmxtools.RequestFromContext(r.Context())
		// log and continue
		slog.Info("handling request", "http.path", r.URL.Path, "http.host", r.Host, "http.method", r.Method, "http.client", r.Header.Get("User-Agent"), "content-type", r.Header.Get("content-type"), "htmxrequest", htmxreq)
		next.ServeHTTP(w, r)
	}
}
