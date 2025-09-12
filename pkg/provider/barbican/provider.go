package barbican

import (
	"context"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/external-secrets/external-secrets/pkg/utils/resolvers"
)

const (
	errGeneric = "barbican provider error: %w"
)

var _ esv1.Provider = &Provider{}

type Provider struct {
	client *gophercloud.ServiceClient
}

func (p *Provider) NewClient(ctx context.Context, store esv1.GenericStore, kube client.Client, namespace string) (esv1.SecretsClient, error) {
	return newClient(ctx, store, kube, namespace)
}

func getProvider(store esv1.GenericStore) (*esv1.BarbicanProvider, error) {
	spec := store.GetSpec()
	if spec.Provider == nil || spec.Provider.Barbican == nil {
		return nil, fmt.Errorf(errGeneric, errors.New("barbican is nil"))
	}
	return spec.Provider.Barbican, nil
}

func (p *Provider) ValidateStore(store esv1.GenericStore) (admission.Warnings, error) {
	if store == nil {
		return nil, fmt.Errorf(errGeneric, errors.New("store is nil"))
	}
	return nil, nil
}

func (p *Provider) Capabilities() esv1.SecretStoreCapabilities {
	return esv1.SecretStoreReadOnly
}

func (p *Provider) Close(ctx context.Context) error {
	return nil
}

func (p *Provider) DeleteSecret(ctx context.Context, ref esv1.PushSecretRemoteRef) error {
	return fmt.Errorf("barbican provider does not support deleting secrets")
}

func (p *Provider) GetAllSecrets(ctx context.Context, ref esv1.ExternalSecretFind) (map[string][]byte, error) {
	return nil, fmt.Errorf("barbican provider does not support get all secrets")
}

func (p *Provider) GetSecret(ctx context.Context, ref esv1.ExternalSecretDataRemoteRef) ([]byte, error) {
	// Implementation for fetching a secret from Barbican would go here.
	return nil, fmt.Errorf("GetSecret method not implemented")
}

func (p *Provider) GetSecretMap(ctx context.Context, ref esv1.ExternalSecretDataRemoteRef) (map[string][]byte, error) {
	return nil, fmt.Errorf("barbican provider does not support get secret map")
}

func (p *Provider) PushSecret(ctx context.Context, secret *corev1.Secret, data esv1.PushSecretData) error {
	return fmt.Errorf("barbican provider does not support pushing secrets")
}

func (p *Provider) SecretExists(ctx context.Context, ref esv1.PushSecretRemoteRef) (bool, error) {
	return false, fmt.Errorf("barbican provider does not support checking if a secret exists")
}

func (p *Provider) Validate() (esv1.ValidationResult, error) {
	return esv1.ValidationResultReady, nil
}

func newClient(ctx context.Context, store esv1.GenericStore, kube client.Client, namespace string) (esv1.SecretsClient, error) {
	provider, err := getProvider(store)
	if err != nil {
		return nil, err
	}
	if provider.AuthURL != "" {
		return nil, fmt.Errorf(errGeneric, errors.New("authURL is required"))
	}

	username, err := resolvers.SecretKeyRef(ctx, kube, store.GetKind(), namespace, provider.Username.SecretRef)
	if err != nil {
		return nil, fmt.Errorf(errGeneric, err)
	}

	password, err := resolvers.SecretKeyRef(ctx, kube, store.GetKind(), namespace, provider.Password.SecretRef)
  if err != nil {
		return nil, fmt.Errorf(errGeneric, err)
	}

	fmt.Println("Resolved username: %s", username)
	fmt.Println("Resolved password: %s", password)

	authopts := gophercloud.AuthOptions{
		IdentityEndpoint: provider.AuthURL,
		TenantName:     provider.TenantName,
		DomainName:     provider.DomainName,
		Username:       username,
		Password:       password,
	}

	auth, err := openstack.AuthenticatedClient(context.TODO(), authopts)
  if err != nil {
    return nil, fmt.Errorf(errGeneric, errors.New("failed to authenticate to OpenStack"))
  }

	  barbicanClient, err := openstack.NewKeyManagerV1(auth, gophercloud.EndpointOpts{
    Region: provider.Region,
  })
	if err != nil {
		return nil, fmt.Errorf(errGeneric, errors.New("failed to create Barbican client"))
	}

	return &Provider{
		client: barbicanClient,
	}, nil
}

func init() {
	esv1.Register(&Provider{}, &esv1.SecretStoreProvider{
		Barbican: &esv1.BarbicanProvider{},
	}, esv1.MaintenanceStatusMaintained)
}
