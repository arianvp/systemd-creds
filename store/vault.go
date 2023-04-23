package store

import (
	"context"
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/approle"
)

type Vault struct {
	ApproleRoleID       string
	ApproleSecretIDFile string
	client              vault.Client
	watcher             *vault.LifetimeWatcher
}

func (v *Vault) Get(ctx context.Context, unitName, credID string) (string, error) {
	kvSecret, err := v.client.KVv2("secret").Get(ctx, unitName)
	if err != nil {
		return "", err
	}

	secret, ok := kvSecret.Data[credID].(string)
	if !ok {
		return "", fmt.Errorf("expected a string value but got %v", kvSecret.Data[credID])
	}
	return secret, nil
}

func (v *Vault) Stop() {
	if v.watcher != nil {
		v.watcher.Stop()
	}
}

func NewVault() (*Vault, error) {
	client, err := vault.NewClient(nil)
	if err != nil {
		return nil, err
	}
	return &Vault{
		ApproleRoleID:       os.Getenv("VAULT_APPROLE_ROLE_ID"),
		ApproleSecretIDFile: os.Getenv("VAULT_APPROLE_SECRET_ID_FILE"),
		client:              *client,
		watcher:             &vault.LifetimeWatcher{},
	}, nil
}

func (v *Vault) Login(ctx context.Context) error {
	approleAuth, err := approle.NewAppRoleAuth(
		v.ApproleRoleID,
		&approle.SecretID{FromFile: v.ApproleSecretIDFile},
		approle.WithWrappingToken(),
	)
	if err != nil {
		return err
	}
	secret, err := v.client.Auth().Login(ctx, approleAuth)
	if err != nil {
		return err
	}
	// TODO: Watchdog needs to take this into account
	watcher, err := v.client.NewLifetimeWatcher(&vault.LifetimeWatcherInput{
		Secret: secret,
	})
	if err != nil {
		return err
	}
	v.watcher = watcher
	go watcher.Start()
	return nil
}
