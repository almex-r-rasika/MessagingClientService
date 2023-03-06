package data

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

/* function for wait till the time
   @param --> ctx, until
   @param value --> Context Object, Time Object
   description --> helper function for wait till the time
   @return --> null
*/
func waitUntil(ctx context.Context, until time.Time) {
	timer := time.NewTimer(time.Until(until))
	defer timer.Stop()

	select {
	case <-timer.C:
		return
	case <-ctx.Done():
		return
	}
}



