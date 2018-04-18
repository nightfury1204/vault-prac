package pkg

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Options struct {
	//cert issuer address
	CertIssuerAddr string

	// vault server address
	VaultAddr string

	// token that has privilege to write certificates
	VaultToken string

	// path in vault where certificates are stored or will be stored
	VaultCertPath string

	// organization name for certificates
	Org string
	// common name for certificates
	Cn string
}

func NewOptions() Options {
	return Options{
		VaultAddr: os.Getenv("VAULT_ADDR"),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.CertIssuerAddr, "cert-issuer-addr", o.CertIssuerAddr, "cert-issuer server address [host:port]")
	fs.StringVar(&o.VaultToken, "token", o.VaultToken, "token that has privilege to write certificates")
	fs.StringVar(&o.VaultAddr, "vault-addr", o.VaultAddr, "vault server address")
	fs.StringVar(&o.VaultCertPath, "vault-cert-path", o.VaultCertPath, "path in vault where certificates are stored or will be stored")
	fs.StringVar(&o.Org, "org", "foo", "organization name for certificates")
	fs.StringVar(&o.Cn, "cn", "bar", "common name for certificates")
}

func (o *Options) Validate() []error {
	var errs []error
	if len(o.VaultAddr) == 0 {
		errs = append(errs, errors.New("vault server address must be non empty"))
	}
	if len(o.CertIssuerAddr) == 0 {
		errs = append(errs, errors.New("cert issuer address must be non empty"))
	}
	if len(o.VaultCertPath) == 0 {
		errs = append(errs, errors.New("path in vault where certificates are stored or will be stored must be non empty"))
	}
	if len(o.VaultToken) == 0 {
		errs = append(errs, errors.New("token must be non empty"))
	}

	return errs
}
