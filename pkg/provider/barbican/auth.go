package barbican

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack"
)

func (c *Client) Authenticate(ctx context.Context) error {
	provider, err := openstack.AuthenticatedClient(c.opts)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	c.provider = provider
	return nil
}

