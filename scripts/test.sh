#!/bin/sh
./scripts/vault.sh
cp ./config/systemd-creds.service /run/systemd/system
cp ./config/systemd-creds.socket /run/systemd/system
systemctl daemon-reload
systemctl start systemd-creds.socket
systemd-run --unit foo  -P -p LoadCredential=bar:/run/systemd-creds/socket sh -c '/run/current-system/sw/bin/cat $CREDENTIALS_DIRECTORY/bar'