// Package alerts reads Azure Monitor metric alert rules for a subscription.
//
// This is a minimal, self-contained demo target for the "Licence to Patch" agent.
// The call below depends on the Azure Monitor metric-alerts REST api-version that
// the armmonitor SDK bakes into every request. Bumping armmonitor v0.12.0 -> v0.13.0
// silently changes that api-version from 2024-03-01-preview to 2026-01-01, which ARM
// rejects with 404 InvalidResourceType at runtime -- while unit tests stay green.
package alerts

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
)

// ListAlertNames returns the names of every metric alert rule in the subscription.
// The api-version used for these calls lives inside the armmonitor module, not here.
func ListAlertNames(ctx context.Context, subscriptionID string, cred azcore.TokenCredential, opts *arm.ClientOptions) ([]string, error) {
	client, err := armmonitor.NewMetricAlertsClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	var names []string
	pager := client.NewListBySubscriptionPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, rule := range page.Value {
			if rule != nil && rule.Name != nil {
				names = append(names, *rule.Name)
			}
		}
	}
	return names, nil
}
