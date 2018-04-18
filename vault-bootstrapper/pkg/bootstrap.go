package pkg

import (
	"io/ioutil"

	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
	"os"
)

const (
	CertStorePathv1 = "secret/certs"

	// https://www.vaultproject.io/api/secret/kv/kv-v2.html
	CertStorePathv2 = "secret/data/certs"

	readOnlyTokenPolicy = `
# Read only permit
path "secret/certs" {
  capabilities = ["read"]
}
`
)

var (
	ttl = os.Getenv("VAULT_TTL")
)

// setting up vault for cert issuer
// task:
//   - if vault is seal, then unseal it
//   - store ca cert and ca key in vault
//   - generate a read access token for cert issuer
func (o *Options) Bootstrap() error {
	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		return errors.Wrap(err, "unable create vault client")
	}

	client.SetAddress(o.Addr)
	client.SetToken(o.VaultToken)

	cert, err := ioutil.ReadFile(o.CaCertFile)
	if err != nil {
		return errors.Wrap(err, "unable to read ca cert read file")
	}

	key, err := ioutil.ReadFile(o.CaKeyFile)
	if err != nil {
		return errors.Wrap(err, "unable to read ca key file")
	}

	data := map[string]interface{}{
		"ca.crt": cert,
		"ca.key": key,
	}

	if client.Sys() == nil {
		return errors.New("nil pointer")
	}

	status, err := client.Sys().SealStatus()
	if err != nil {
		return errors.Wrap(err, "unable to get vault seal status")
	}

	glog.Infoln("---------------------------------------------")
	glog.Infoln("vault seal status:")
	glog.Infoln(status)
	glog.Infoln("---------------------------------------------")

	// if vault is sealed then unseal it
	// if vault is in dev mode then ignor it
	if !o.IsDev && status.Sealed {
		if status.T > len(o.UnSealKeys) {
			return errors.Errorf("insufficient number of unseal keys. %d keys required", status.T)
		}

		unseal := false
		for i := range o.UnSealKeys {
			resp, err := client.Sys().Unseal(o.UnSealKeys[i])
			if err != nil {
				return errors.Wrap(err, "error occurred when unsealing vault")
			}

			if !resp.Sealed {
				unseal = true
				break
			}
		}

		if !unseal {
			return errors.New("unsealed failed")
		}
	}

	// writing ca crt/key in vault
	_, err = client.Logical().Write(CertStorePathv1, data)
	if err != nil {
		return errors.Wrap(err, "unable write data in vault path "+CertStorePathv1)
	}

	// creating read only policy for cert issuer

	glog.Infoln("---------------------------------------------")
	glog.Infoln("readCA policy:")
	glog.Infoln(readOnlyTokenPolicy)
	glog.Infoln("---------------------------------------------")

	policyName := "readCA"
	err = client.Sys().PutPolicy(policyName, readOnlyTokenPolicy)
	if err != nil {
		return errors.Wrap(err, "unable to create policy")
	}

	// generating token for cert issuer
	tokenReq := &vaultapi.TokenCreateRequest{
		Policies: []string{policyName},
		Metadata: map[string]string{
			"user": "cert-issuer",
		},
		DisplayName: "Read-CA-Token",
		NoParent:    true,
		Period:      ttl,
		TTL:         ttl,
	}
	secret, err := client.Auth().Token().Create(tokenReq)
	if err != nil {
		return errors.Wrap(err, "unable to create token for cert issuer")
	}

	glog.Infoln("--------------------------------")
	glog.Infoln("Cert Issuer token")
	value, _ := json.MarshalIndent(secret, "", "  ")
	glog.Infoln(string(value))
	glog.Infoln("--------------------------------")

	fmt.Println("Token for cert-issuer : "+secret.Auth.ClientToken)

	return nil
}
