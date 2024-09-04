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
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"

	log "github.com/golang/glog"
)

const (
	// cookieName is the name used to store authen details.
	// Those details are going to be:
	// concat(username, "::", time.Now().Unix(), "::",  sha512(username, time.Now().Unix(), masterSecret))
	// The cookie should auto-expire, but if the name/time don't match the hash expire directly.
	cookieName   = "shizzle"
	masterSecret = "508be247d847fe6061c4f912c7a0d9717daa76f19b0690534295c29796294550c95884f1ccef8a6667049728b4e10e91cc0ab6f831aa58eafc25c5d64474c377"
)

var (
	addr  = flag.String("addr", "127.0.0.1", "Address to listen for connections.")
	port  = flag.Int("port", 9876, "Port on which to listen for this service.")
	userN = flag.String("un", "bob", "User to use as a login to the web application.")
	userP = flag.String("pw", "bobs-pass", "User's password for the web application.")

	// Default maxAge for a cookie.
	cookieMaxAge = 86400 * 14 // 14 days before expiry of the cookie.
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

func (h *handler) calcValue(r *http.Request) string {
	// Method is already POST only by here, so just extract the POST data.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Infof("failed to read body content: %v", err)
		return ""
	}
	log.Infof("BODY: %v", body)
	return "fee"
}

func (h *handler) updateCookie(c *http.Cookie) *http.Cookie {
	return c
}

// Either update an existing cookie, or create a new cookie.
func (h *handler) createCookie(r *http.Request) *http.Cookie {
	// Retrive the current cookie(s),
	c, err := r.Cookie(cookieName)
	if err != nil {
		// A nil reply means the cookie is not there, create it.
		c = &http.Cookie{
			Name:     cookieName,
			Value:    h.calcValue(r),
			Path:     "/",
			MaxAge:   cookieMaxAge,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
		return c
	}
	c = h.updateCookie(c)
	return c
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
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, d.String())
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Infof("failed to read body content: %v", err)
			fmt.Fprintf(w, d.String())
			break
		}
		log.Infof("BODY: %v", body)
		fmt.Fprintf(w, d.String())
	}
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
