# senerio 1

```
$ vault server -dev

```

Open another terminal and run following command 
```
$ export VAULT_ADDR='http://127.0.0.1:8200'

```


Using K/v version 1 (unversioned key value pair)

```
# by default secret/ is versionsed K/V
$ vault secrets disable secret/
Success! Disabled the secrets engine (if it existed) at: secret/

$ vault secrets enable -version=1 -path=secret kv
Success! Enabled the kv secrets engine at: secret/
```

For K/V version 2 (versioned key value pair) [link](https://www.vaultproject.io/api/secret/kv/kv-v2.html)

Run vault-bootstrapper
```
vault-bootstrappers run \
   --dev \
   --token=3a88f965-5bbb-b5a5-bca6-7445b9d41eb7 \
   --ca-cert-file=/home/ac/go/src/github.com/nightfury1204/vault-prac/vault-bootstrapper/dist/certs/ca.crt \
   --ca-key-file=/home/ac/go/src/github.com/nightfury1204/vault-prac/vault-bootstrapper/dist/certs/ca.key
```

