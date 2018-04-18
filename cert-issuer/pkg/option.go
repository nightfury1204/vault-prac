package pkg

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Options struct {
	Port string
	// vault server address
	VaultAddr string

	// token that has privilege to write certificates
	VaultToken string
}

func NewOptions() Options {
	return Options{
		VaultAddr: os.Getenv("VAULT_ADDR"),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Port, "port", "10088", "cert-issuer server port")
	fs.StringVar(&o.VaultToken, "token", o.VaultToken, "token that has privilege to read certificates")
	fs.StringVar(&o.VaultAddr, "vault-addr", o.VaultAddr, "vault server address")
}

func (o *Options) Validate() []error {
	var errs []error
	if len(o.VaultAddr) == 0 {
		errs = append(errs, errors.New("vault server address must be non empty"))
	}
	if len(o.VaultToken) == 0 {
		errs = append(errs, errors.New("token must be non empty"))
	}

	return errs
}
