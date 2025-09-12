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
	"github.com/external-secrets/external-secrets/pkg/utils/resolvers"
)

const (
	errGeneric = "barbican provider error: %w"
)

var _ esv1.Provider = &Provider{}

type Provider struct {}

func (p *Provider) Capabilities() esv1.SecretStoreCapabilities {
	return esv1.SecretStoreReadOnly
}

func (p *Provider) ValidateStore(store esv1.GenericStore) (admission.Warnings, error) {
	if store == nil {
		return nil, fmt.Errorf(errGeneric, errors.New("store is nil"))
	}
	return nil, nil
}

func (p *Provider) NewClient(ctx context.Context, store esv1.GenericStore, kube client.Client, namespace string) (esv1.SecretsClient, error) {
	return newClient(ctx, store, kube, namespace)
}

func getProvider(store esv1.GenericStore) (*esv1.BarbicanProvider, error) {
	spec := store.GetSpec()
	if spec.Provider == nil || spec.Provider.Barbican == nil {
		return nil, fmt.Errorf(errGeneric, errors.New("barbican is nil"))
	}
	fmt.Printf("Barbican provider config: %+v\n", spec.Provider.Barbican)
	return spec.Provider.Barbican, nil
}

func newClient(ctx context.Context, store esv1.GenericStore, kube client.Client, namespace string) (esv1.SecretsClient, error) {
	provider, err := getProvider(store)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Using Barbican provider config: %+v\n", provider)
	if provider.AuthURL == "" {
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

	c := &Client{
		  keyManager: barbicanClient,
	}

	return c, nil
}

func init() {
	esv1.Register(&Provider{}, &esv1.SecretStoreProvider{
		Barbican: &esv1.BarbicanProvider{},
	}, esv1.MaintenanceStatusMaintained)
}
