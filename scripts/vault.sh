#!/bin/sh
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=lol
vault auth enable approle || true
vault write -force auth/approle/role/utm 
vault read auth/approle/role/utm/role-id
vault write -wrap-ttl=120s -field wrapping_token -f auth/approle/role/utm/secret-id > /run/credstore/approle-secret-id
gg
vault kv put -mount=secret foo.service bar=baz
vault policy write utm - <<EOF
path "secret/data/foo.service" {
  capabilities = ["read"]
}
EOF
vault write auth/approle/role/utm token_policies=utm