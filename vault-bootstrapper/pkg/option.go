package pkg

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Options struct {
	// vault server address
	Addr string

	// is vault server running in dev mode
	IsDev bool

	CaCertFile string

	CaKeyFile string

	// token that has privilege to write certificates
	VaultToken string

	// Keys that requires to unseal vault
	UnSealKeys []string
}

func NewOptions() Options {
	return Options{
		Addr: os.Getenv("VAULT_ADDR"),
	}
}

func (o *Options) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.CaCertFile, "ca-cert-file", o.CaCertFile, "CA cert file path")
	fs.StringVar(&o.CaKeyFile, "ca-key-file", o.CaCertFile, "CA key file path")
	fs.BoolVar(&o.IsDev, "dev", o.IsDev, "Is vault server running in dev mode")
	fs.StringSliceVar(&o.UnSealKeys, "unseal-keys", o.UnSealKeys, "Keys that requires to unseal vault")
	fs.StringVar(&o.VaultToken, "token", o.VaultToken, "token that has privilege to write certificates")
	fs.StringVar(&o.Addr, "vault-addr", o.Addr, "vault server address")
}

func (o *Options) Validate() []error {
	var errs []error
	if len(o.Addr) == 0 {
		errs = append(errs, errors.New("vault server address must be non empty"))
	}
	if len(o.CaCertFile) == 0 {
		errs = append(errs, errors.New("ca cert file path must be non empty"))
	}
	if len(o.CaKeyFile) == 0 {
		errs = append(errs, errors.New("ca key file path must be non empty"))
	}
	if o.IsDev == false {
		if len(o.UnSealKeys) == 0 {
			errs = append(errs, errors.New("unseal keys required to unseal vault"))
		}
	}
	if len(o.VaultToken) == 0 {
		errs = append(errs, errors.New("token must be non empty"))
	}

	return errs
}
