package barbican

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/gophercloud/gophercloud/v2/openstack"
  "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/secrets"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
)

var _ esv1.SecretsClient = &Client{}

func (c *Client) GetAllSecrets() {


}

func (c *Client) GetSecret() {

}

func (c *Client) GetSecretMap() {

}

func (c *Client) Validate() {

}

func (c *Client) SecretExists() {

}
