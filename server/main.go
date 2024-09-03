// Simple http server application with three endpoints:
//
//	/login - login to the server, set an auth cookie.
//	/stuff - show a page only after authentication.
//	/logout - destroy the auth cookie.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	log "github.com/golang/glog"
)

var (
	addr  = flag.String("addr", "127.0.0.1", "Address to listen for connections.")
	port  = flag.Int("port", 9876, "Port on which to listen for this service.")
	userN = flag.String("un", "bob", "User to use as a login to the web application.")
	userP = flag.String("pw", "bobs-pass", "User's password for the web application.")
)

type handler struct {
	user, passwd string
	tmpls        map[string]*template.Template
}

func newHandler(u, p string) (*handler, error) {
	tmpls := map[string]*template.Template{
		"index":  indexTmpl,
		"login":  loginTmpl,
		"logout": logoutTmpl,
		"stuff":  stuffTmpl,
	}
	return &handler{
		user:   u,
		passwd: p,
		tmpls:  tmpls,
	}, nil
}

func (h *handler) index(w http.ResponseWriter, r *http.Request) {
	var d bytes.Buffer
	err := h.tmpls["index"].Execute(&d, nil)
	if err != nil {
		fmt.Fprintf(w, "failed to exec index template: %v", err)
	}
	fmt.Fprintf(w, d.String())
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	var d bytes.Buffer
	err := h.tmpls["login"].Execute(&d, nil)
	if err != nil {
		fmt.Fprintf(w, "failed to exec login template: %v", err)
	}
	fmt.Fprintf(w, d.String())
}

func (h *handler) logout(w http.ResponseWriter, r *http.Request) {
	var d bytes.Buffer
	err := h.tmpls["logout"].Execute(&d, nil)
	if err != nil {
		fmt.Fprintf(w, "failed to exec logout template: %v", err)
	}
	fmt.Fprintf(w, d.String())
}

func (h *handler) stuff(w http.ResponseWriter, r *http.Request) {
	var d bytes.Buffer
	err := h.tmpls["stuff"].Execute(&d, nil)
	if err != nil {
		fmt.Fprintf(w, "failed to exec stuff template: %v", err)
	}
	fmt.Fprintf(w, d.String())
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/login"):
		log.Info("Got a login request.")
		h.login(w, r)
	case strings.HasPrefix(r.URL.Path, "/logout"):
		log.Info("Got a logout request.")
		h.logout(w, r)
	case strings.HasPrefix(r.URL.Path, "/stuff"):
		log.Info("Got a stuff request.")
		h.stuff(w, r)
	default:
		log.Info("Got a index request.")
		h.index(w, r)
	}
}

func main() {
	flag.Parse()
	worker, err := newHandler(*userN, *userP)
	if err != nil {
		fmt.Printf("failed to create handler: %v", err)
		return
	}

	h := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", *addr, *port),
		Handler: worker,
	}
	log.Fatal(h.ListenAndServe())
}
