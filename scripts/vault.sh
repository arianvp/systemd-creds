#!/bin/sh
export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=lol
vault auth enable approle || true
vault write -force auth/approle/role/utm role_id=462cd965-ecf1-6ce2-ae50-f91bde23e001
vault read auth/approle/role/utm/role-id
vault write -wrap-ttl=120s -field wrapping_token -f auth/approle/role/utm/secret-id > /run/credstore/approle-secret-id
vault kv put -mount=secret foo.service bar=baz
vault policy write utm - <<EOF
path "secret/data/foo.service" {
  capabilities = ["read"]
}
EOF
vault write auth/approle/role/utm token_policies=utm
cp ./config/systemd-creds.service /run/systemd/system
cp ./config/systemd-creds.socket /run/systemd/system
systemctl daemon-reload
systemctl start systemd-creds.socket
systemd-run --unit foo  -P -p LoadCredential=bar:/run/systemd-creds/socket sh -c '/run/current-system/sw/bin/cat $CREDENTIALS_DIRECTORY/bar'