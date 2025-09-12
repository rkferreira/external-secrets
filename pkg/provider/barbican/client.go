package barbican

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"

	corev1 "k8s.io/api/core/v1"

	esapi "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
)

const (
	errClientGeneric = "barbican provider: %w"
)

var _ esapi.SecretsClient = &Client{}

type Client struct {
	keyManager *gophercloud.ServiceClient
}

func (c *Client) GetAllSecrets(ctx context.Context, ref esapi.ExternalSecretFind) (map[string][]byte, error) {
	return nil, fmt.Errorf("barbican provider does not support get all secrets")
}

func (c *Client) GetSecret(ctx context.Context, ref esapi.ExternalSecretDataRemoteRef) ([]byte, error) {
	// Implementation for fetching a secret from Barbican would go here.
	return nil, fmt.Errorf("GetSecret method not implemented")
}

func (c *Client) GetSecretMap(ctx context.Context, ref esapi.ExternalSecretDataRemoteRef) (map[string][]byte, error) {
	return nil, fmt.Errorf("barbican provider does not support get secret map")
}

func (c *Client) PushSecret(ctx context.Context, secret *corev1.Secret, data esapi.PushSecretData) error {
	return fmt.Errorf("barbican provider does not support pushing secrets")
}

func (c *Client) SecretExists(ctx context.Context, ref esapi.PushSecretRemoteRef) (bool, error) {
	return false, fmt.Errorf("barbican provider does not support checking if a secret exists")
}

func (c *Client) DeleteSecret(ctx context.Context, ref esapi.PushSecretRemoteRef) error {
	return fmt.Errorf("barbican provider does not support deleting secrets")
}

func (c *Client) Validate() (esapi.ValidationResult, error) {
	return esapi.ValidationResultReady, nil
}

func (c *Client) Close(ctx context.Context) error {
	return nil
}
