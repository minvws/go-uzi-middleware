package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"

	uzi "github.com/minvws/go-uzimiddleware"
)

func main() {
	uziMiddleware := uzi.New(uzi.Options{
		StrictCACheck: false,
		AllowedTypes:  []uzi.UziType{uzi.UziTypeCareProvider, uzi.UziTypeNamedEmployee},
		AllowedRoles:  []uzi.UziRole{uzi.UziRoleDoctor},
		Debug:         false,
	})

	server := &http.Server{
		TLSConfig: createTLSConfig("./uzi.client_ca.cert"),
		Addr:      ":8041",
		Handler:   uziMiddleware.Handler(http.HandlerFunc(handler)),
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

func handler(w http.ResponseWriter, r *http.Request) {
	// Get UZI user from context
	uziUser := r.Context().Value(uzi.UziContext("uzi")).(*uzi.UziUser)

	w.Header().Set("Content-Type", "text/html")
	_, _ = w.Write([]byte("Logged in: " + uziUser.SurName))
}
