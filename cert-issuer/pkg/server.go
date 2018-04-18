package pkg

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/bmizerany/pat"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/cert"
	"encoding/base64"
)

const (
	CertStorePathv1 = "secret/certs"
)

type certStore struct {
	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
}

func New() *certStore {
	return &certStore{}
}

// issue client crt/key
//
// query parameter: org, cn
// here, cn=common name, org=organization name
// example: 127.0.0.1:10088/issue/cert?org=foo&cn=bar
//
// response format:
// {
//    "crt" : "...base64encoded..cert.."
//	  "key" : "...base64encoded..key..."
// }
func (c *certStore) IssueCert(w http.ResponseWriter, r *http.Request) {
	org, found := r.URL.Query()["org"]
	if !found {
		writeErrorResponse(w, "empty org name")
		return
	}

	cn, found := r.URL.Query()["cn"]
	if !found || len(cn)!=1 {
		writeErrorResponse(w, "empty/multiple cn name")
		return
	}
	cfg := cert.Config{
		Organization: org,
		CommonName: cn[0],
		Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	key, err := cert.NewPrivateKey()
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "failed to generate private key").Error())
		return
	}
	crt, err := cert.NewSignedCert(cfg, key, c.caCert, c.caKey)
	if err != nil {
		writeErrorResponse(w, errors.Wrap(err, "failed to generate server certificate").Error())
		return
	}

	// base64 encode
	encodedCrt := base64.StdEncoding.EncodeToString(cert.EncodeCertPEM(crt))
	encodedKey := base64.StdEncoding.EncodeToString(cert.EncodePrivateKeyPEM(key))

	writeResponse(w,encodedCrt, encodedKey)
	return
}

func (o *Options) GetCA() (*certStore, error) {
	store := New()

	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return nil, errors.Wrap(err, "unable create vault client")
	}

	client.SetAddress(o.VaultAddr)
	client.SetToken(o.VaultToken)

	data, err := client.Logical().Read(CertStorePathv1)
	if err != nil {
		return nil, errors.Wrap(err, "unable read ca cert/key from vault path "+CertStorePathv1)
	}

	if caCert, found := data.Data["ca.crt"]; found {
		crt, err := cert.ParseCertsPEM(intefaceToBytes(caCert))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse ca cert")
		}
		store.caCert = crt[0]
	} else {
		return nil, errors.New("ca.crt not found")
	}

	if caKey, found := data.Data["ca.key"]; found {
		key, err := cert.ParsePrivateKeyPEM(intefaceToBytes(caKey))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse ca key")
		}
		store.caKey = key.(*rsa.PrivateKey)
	} else {
		return nil, errors.New("ca.key not found")
	}

	return store, nil
}

func (o *Options) RunServer() error {

	store, err := o.GetCA()
	if err != nil {
		return errors.Wrap(err, "unable to get CA crt/key")
	}

	m := pat.New()

	m.Get("/issue/cert", http.HandlerFunc(store.IssueCert))

	server := &http.Server{
		Addr:         "127.0.0.1:" + o.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
	}

	fmt.Println("Cert issuer server address: " + server.Addr)

	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func writeErrorResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{ "error" : "%s"}`, message)))
}

func writeResponse(w http.ResponseWriter, crt string, key string) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{ "crt" : "%s", "key" : "%s"}`, crt, key)))
}

func intefaceToBytes(v interface{}) []byte {
	switch v.(type) {
	case string:
		data, _ := base64.StdEncoding.DecodeString(v.(string))
		return data
	}
	return nil
}
