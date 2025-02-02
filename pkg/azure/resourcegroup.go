package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/azure/spin-aks-plugin/pkg/logger"
)

func ListResourceGroups(ctx context.Context, sub string) ([]armresources.ResourceGroup, error) {
	lgr := logger.FromContext(ctx).With("subscription", sub)
	ctx = logger.WithContext(ctx, lgr)
	lgr.Debug("listing Azure resource groups")

	cred, err := getCred()
	if err != nil {
		return nil, fmt.Errorf("getting credentials: %w", err)
	}

	client, err := armresources.NewResourceGroupsClient(sub, cred, nil)
	if err != nil {
		return nil, fmt.Errorf("creating resource groups client: %w", err)
	}

	var rgs []armresources.ResourceGroup
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing resource groups page: %w", err)
		}

		for _, rg := range page.Value {
			if rg == nil {
				return nil, errors.New("nil rg")
			}

			rgs = append(rgs, *rg)
		}
	}

	lgr.Debug("finished listing Azure resource groups")
	return rgs, nil
}
