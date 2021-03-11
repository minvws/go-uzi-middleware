package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uzi "github.com/minvws/go-uzimiddleware"
	"github.com/urfave/negroni"
)

func main() {
	uziMiddleware := uzi.New(uzi.Options{
		StrictCACheck: false,
		AllowedTypes:  []uzi.UziType{uzi.UZI_TYPE_CARE_PROVIDER, uzi.UZI_TYPE_NAMED_EMPLOYEE},
		AllowedRoles:  []uzi.UziRole{uzi.UZI_ROLE_DOCTOR},
		Debug:         false,
	})

	r := mux.NewRouter()
	r.HandleFunc("/", defaultHandler)
	r.Handle("/uzi", negroni.New(
		negroni.HandlerFunc(uziMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(secureHandler)),
	))

	server := &http.Server{
		TLSConfig: createTLSConfig("./uzi.client_ca.cert"),
		Addr:      ":8041",
		Handler:   r,
	}
	server.ListenAndServeTLS("./server.crt", "./server.key")
}

func createTLSConfig(certPath string) *tls.Config {
	certs := x509.NewCertPool()
	pemData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Fatal("cannot read certificate: ", err)
	}
	certs.AppendCertsFromPEM(pemData)

	return &tls.Config{
		ClientAuth:         tls.RequestClientCert,
		RootCAs:            certs,
		InsecureSkipVerify: false,
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte("Public available"))
}

func secureHandler(w http.ResponseWriter, r *http.Request) {
	// Get UZI user from context
	uziUser := r.Context().Value("uzi").(*uzi.UziUser)

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte("Logged in: " + uziUser.SurName))
}
