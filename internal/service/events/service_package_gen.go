// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package events

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory:  newEventBusesDataSource,
			TypeName: "aws_cloudwatch_event_buses",
			Name:     "Event Buses",
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceBus,
			TypeName: "aws_cloudwatch_event_bus",
			Name:     "Event Bus",
		},
		{
			Factory:  dataSourceConnection,
			TypeName: "aws_cloudwatch_event_connection",
			Name:     "Connection",
		},
		{
			Factory:  dataSourceSource,
			TypeName: "aws_cloudwatch_event_source",
			Name:     "Source",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceAPIDestination,
			TypeName: "aws_cloudwatch_event_api_destination",
			Name:     "API Destination",
		},
		{
			Factory:  resourceArchive,
			TypeName: "aws_cloudwatch_event_archive",
			Name:     "Archive",
		},
		{
			Factory:  resourceBus,
			TypeName: "aws_cloudwatch_event_bus",
			Name:     "Event Bus",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceBusPolicy,
			TypeName: "aws_cloudwatch_event_bus_policy",
			Name:     "Event Bus Policy",
		},
		{
			Factory:  resourceConnection,
			TypeName: "aws_cloudwatch_event_connection",
			Name:     "Connection",
		},
		{
			Factory:  resourceEndpoint,
			TypeName: "aws_cloudwatch_event_endpoint",
			Name:     "Global Endpoint",
		},
		{
			Factory:  resourcePermission,
			TypeName: "aws_cloudwatch_event_permission",
			Name:     "Permission",
		},
		{
			Factory:  resourceRule,
			TypeName: "aws_cloudwatch_event_rule",
			Name:     "Rule",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceTarget,
			TypeName: "aws_cloudwatch_event_target",
			Name:     "Target",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Events
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*eventbridge.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))

	return eventbridge.NewFromConfig(cfg,
		eventbridge.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
	), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
