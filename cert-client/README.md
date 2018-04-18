# Cert Client

- Check crt/key in vault
- Request cert issuer for crt/key
- store it in vault

```
$ ./cert-client run -h
get client cert

Usage:
  cert-client run [flags]

Flags:
      --cert-issuer-addr string   cert-issuer server address [host:port]
      --cn string                 common name for certificates (default "bar")
  -h, --help                      help for run
      --org string                organization name for certificates (default "foo")
      --token string              token that has privilege to write certificates
      --vault-addr string         vault server address
      --vault-cert-path string    path in vault where certificates are stored or will be stored

Global Flags:
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging

```

Provided vault token to cert client must have read and write access to vault cert path

Keys used for storing crt/key
```
key     value
-----  -------
crt    <base64encoded>
key    <base64encoded>
```