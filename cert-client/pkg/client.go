package pkg

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"

	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/cert"
)

type certStore struct {
	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
}

func New() *certStore {
	return &certStore{}
}

// store crt/key in vault
//
// key in vault
//    for crt : crt
//    for key : key
//
// parameter crt and key are base64 encoded
func (o *Options) StoreInVault(crt string, key string) error {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return errors.Wrap(err, "unable create vault client")
	}

	client.SetAddress(o.VaultAddr)
	client.SetToken(o.VaultToken)

	data := map[string]interface{}{
		"crt": crt,
		"key": key,
	}

	_, err = client.Logical().Write(o.VaultCertPath, data)
	if err != nil {
		return errors.Wrap(err, "unable store crt/key in vault path "+o.VaultCertPath)
	}
	return nil
}

// To ger crt/key pair
// first check crt/key in the vault storage
// if not found, then request crt/key from cert issuer and store it in vault
//
// key in vault
//    for crt : crt
//    for key : key
func (o *Options) GetCert() (*certStore, error) {
	store := New()

	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return nil, errors.Wrap(err, "unable create vault client")
	}

	client.SetAddress(o.VaultAddr)
	client.SetToken(o.VaultToken)

	data, err := client.Logical().Read(o.VaultCertPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable read ca cert/key from vault path "+o.VaultCertPath)
	}

	crtFound := true
	if data!=nil {
		if caCert, found := data.Data["crt"]; found {
			crt, err := cert.ParseCertsPEM(intefaceToBytes(caCert))
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse ca cert")
			}
			store.caCert = crt[0]
		} else {
			crtFound = false
		}

		if caKey, found := data.Data["key"]; found {
			key, err := cert.ParsePrivateKeyPEM(intefaceToBytes(caKey))
			if err != nil {
				return nil, errors.Wrap(err, "unable to parse key")
			}
			store.caKey = key.(*rsa.PrivateKey)
		} else {
			crtFound = false
		}
	} else {
		crtFound = false
	}

	if !crtFound {
		// request for crt/key for cert-issuer
		resp, err := http.Get(fmt.Sprintf("http://%s/issue/cert?org=%s&cn=%s", o.CertIssuerAddr, o.Org, o.Cn))
		if err != nil {
			return nil, errors.Wrap(err, "unable get crt/key from cert issuer")
		}

		data := struct {
			Crt   string `json:"crt"`
			Key   string `json:"key"`
			Error string `json:"error"`
		}{}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read response body from cert issuer")
		}

		err = json.Unmarshal(contents, &data)
		if err != nil {
			return nil, errors.Wrap(err, "unable to unmarshall response from cert issuer")
		}

		if data.Error != "" {
			return nil, errors.Errorf("failed to get crt/key. reason: %v", data.Error)
		}

		crt, err := cert.ParseCertsPEM(intefaceToBytes(data.Crt))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse ca cert")
		}
		store.caCert = crt[0]

		key, err := cert.ParsePrivateKeyPEM(intefaceToBytes(data.Key))
		if err != nil {
			return nil, errors.Wrap(err, "unable to parse key")
		}
		store.caKey = key.(*rsa.PrivateKey)

		// store in vault
		err = o.StoreInVault(data.Crt, data.Key)
		if err != nil {
			return nil, errors.Wrap(err, "unable to store crt/key")
		}
	}

	return store, nil
}

func (o *Options) Run() error {

	_, err := o.GetCert()
	if err != nil {
		return errors.Wrap(err, "unable to get client crt/key")
	}

	fmt.Println("Getting crt/key is successful")

	return nil
}

// decode base64 encoded interface
func intefaceToBytes(v interface{}) []byte {
	switch v.(type) {
	case string:
		data, _ := base64.StdEncoding.DecodeString(v.(string))
		return data
	}
	return nil
}
