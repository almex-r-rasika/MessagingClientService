package data

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"
)

/* function for add ca certificate
   @param --> null
   @param value --> null
   description --> add ca certificate to call the api
   @return --> null
*/
func addCaCertificate() http.Client{

	absPath, _ := filepath.Abs("../main-service/data/smapa.crt")

    pem, err := os.ReadFile(absPath)
	if err != nil {
        Log.Fatal(err.Error())
	}

    caCertPool, err := x509.SystemCertPool()
	if err != nil {
		Log.Fatal(err.Error())
	}

	if !caCertPool.AppendCertsFromPEM(pem) {
		Log.Fatal("failed to add ca cert")
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		Log.Fatal("invalid default transport")
	}

	transport := defaultTransport.Clone()

	transport.TLSClientConfig = &tls.Config{
		RootCAs: caCertPool,
	}

    client := http.Client{
		Transport: transport,
	}

	return client
}




