package uzimiddleware

import (
	"context"
	"errors"
	"log"
	"net/http"
)

var errPeerCertNotFound = errors.New("uzi: peer certificate not found")

// errorHandler is a function that is called when an error occurs
type errorHandler func(w http.ResponseWriter, r *http.Request, err string)

// Options allows you to easily set options for the uzi middleware
type Options struct {
	StrictCACheck bool      // Strict check on the CA
	AllowedTypes  []UziType // Allowed card types
	AllowedRoles  []UziRole // Allowed card roles

	ErrorHandler errorHandler // Custom error handler
	Debug        bool         // outputs debug information when set to true
}

// UZIMiddleware is the actual middleware
type UZIMiddleware struct {
	Options Options
}

// DefaultErrorHandler will be used when no custom error handler is set through the options
func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err string) {
	http.Error(w, err, http.StatusUnauthorized)
}

// New initializes a new middleware structure with options.
func New(options Options) *UZIMiddleware {
	if options.ErrorHandler == nil {
		options.ErrorHandler = DefaultErrorHandler
	}

	return &UZIMiddleware{
		Options: options,
	}
}

// HandlerWithNext will authenticate through UZI, and calls the next middleware handler when failed
func (uzi *UZIMiddleware) HandlerWithNext(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := uzi.CheckUziCertificate(w, r)
	if err == nil && next != nil {
		next(w, r)
	}
}

// CheckUziCertificate will verify the UZI information found in the given request
func (uzi *UZIMiddleware) CheckUziCertificate(w http.ResponseWriter, r *http.Request) error {
	if len(r.TLS.PeerCertificates) == 0 {
		err := errPeerCertNotFound

		uzi.logf("uzi: peer certificate not found.")

		uzi.Options.ErrorHandler(w, r, err.Error())
		return err
	}

	cert := r.TLS.PeerCertificates[0]
	user, err := NewUziUserFromCert(cert)
	if err != nil {
		uzi.logf("uzi: cannot create uzi user from certificate: %s", err.Error())

		uzi.Options.ErrorHandler(w, r, err.Error())
		return err
	}

	err = user.Validate(uzi.Options.StrictCACheck, uzi.Options.AllowedTypes, uzi.Options.AllowedRoles)
	if err != nil {
		uzi.logf("uzi: cannot validate uzi user: %s", err)

		uzi.Options.ErrorHandler(w, r, err.Error())
		return err
	}

	// Add the UZI user to the request
	newReq := r.WithContext(context.WithValue(r.Context(), "uzi", user))
	*r = *newReq

	return nil
}

// logf will output log statements when the debug option is enabled
func (m *UZIMiddleware) logf(format string, args ...interface{}) {
	if m.Options.Debug {
		log.Printf(format, args...)
	}
}
