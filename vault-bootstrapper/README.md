# Vault bootstrapper

- Unseal vault if required
- Store Ca crt/key in vault
- Create read access token for cert issuer

```
$ ./vault-bootstrapper run -h
  bootstrapper
  
  Usage:
    vault-bootstrapper run [flags]
  
  Flags:
        --ca-cert-file string   CA cert file path
        --ca-key-file string    CA key file path
        --dev                   Is vault server running in dev mode
    -h, --help                  help for run
        --token string          token that has privilege to write certificates
        --unseal-keys strings   Keys that requires to unseal vault
        --vault-addr string     vault server address
  
  Global Flags:
        --alsologtostderr                  log to standard error as well as files
        --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
        --log_dir string                   If non-empty, write log files in this directory
        --logtostderr                      log to standard error instead of files
        --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
    -v, --v Level                          log level for V logs
        --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
 
```

