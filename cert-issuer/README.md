# senerio 1
Issue client crt/key
```
cert-issuer

Usage:
  cert-issuer run [flags]

Flags:
  -h, --help                help for run
      --port string         cert-issuer server port (default "10088")
      --token string        token that has privilege to write certificates
      --vault-addr string   vault server address

```


Request url:
```
127.0.0.1:10088/issue/cert?org=foo&cn=bar 
```

Response Format:
```json
{
  "crt" : "<base64encoded>",
  "key" : "<base64encoded>"
}
```


