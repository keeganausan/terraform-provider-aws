// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package connect

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceBotAssociation,
			TypeName: "aws_connect_bot_association",
			Name:     "Bot Association",
		},
		{
			Factory:  dataSourceContactFlow,
			TypeName: "aws_connect_contact_flow",
			Name:     "Contact Flow",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceContactFlowModule,
			TypeName: "aws_connect_contact_flow_module",
			Name:     "Contact Flow Module",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceHoursOfOperation,
			TypeName: "aws_connect_hours_of_operation",
			Name:     "Hours Of Operation",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceInstance,
			TypeName: "aws_connect_instance",
			Name:     "Instance",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceInstanceStorageConfig,
			TypeName: "aws_connect_instance_storage_config",
			Name:     "Instance Storage Config",
		},
		{
			Factory:  dataSourceLambdaFunctionAssociation,
			TypeName: "aws_connect_lambda_function_association",
			Name:     "Lambda Function Association",
		},
		{
			Factory:  dataSourcePrompt,
			TypeName: "aws_connect_prompt",
			Name:     "Prompt",
		},
		{
			Factory:  dataSourceQueue,
			TypeName: "aws_connect_queue",
			Name:     "Queue",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceQuickConnect,
			TypeName: "aws_connect_quick_connect",
			Name:     "Quick Connect",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceRoutingProfile,
			TypeName: "aws_connect_routing_profile",
			Name:     "Routing Profile",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceSecurityProfile,
			TypeName: "aws_connect_security_profile",
			Name:     "Security Profile",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  DataSourceUser,
			TypeName: "aws_connect_user",
			Name:     "User",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceUserHierarchyGroup,
			TypeName: "aws_connect_user_hierarchy_group",
			Name:     "User Hierarchy Group",
			Tags:     &types.ServicePackageResourceTags{},
		},
		{
			Factory:  dataSourceUserHierarchyStructure,
			TypeName: "aws_connect_user_hierarchy_structure",
			Name:     "User Hierarchy Structure",
		},
		{
			Factory:  dataSourceVocabulary,
			TypeName: "aws_connect_vocabulary",
			Name:     "Vocabulary",
			Tags:     &types.ServicePackageResourceTags{},
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceAssociatePhoneNumberContactFlow,
			TypeName: "aws_connect_associate_phone_number_contact_flow",
			Name:     "Associate Phone Number Contact Flow",
    },
    {
			Factory:  resourceAssociateQueueQuickConnects,
			TypeName: "aws_connect_associate_queue_quick_connects",
			Name:     "Associate Queue Quick Connects",
		},
		{
			Factory:  resourceBotAssociation,
			TypeName: "aws_connect_bot_association",
			Name:     "Bot Association",
		},
		{
			Factory:  resourceContactFlow,
			TypeName: "aws_connect_contact_flow",
			Name:     "Contact Flow",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceContactFlowModule,
			TypeName: "aws_connect_contact_flow_module",
			Name:     "Contact Flow Module",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceHoursOfOperation,
			TypeName: "aws_connect_hours_of_operation",
			Name:     "Hours Of Operation",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceInstance,
			TypeName: "aws_connect_instance",
			Name:     "Instance",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceInstanceStorageConfig,
			TypeName: "aws_connect_instance_storage_config",
			Name:     "Instance Storage Config",
		},
		{
			Factory:  resourceLambdaFunctionAssociation,
			TypeName: "aws_connect_lambda_function_association",
			Name:     "Lambda Function Association",
		},
		{
			Factory:  resourcePhoneNumber,
			TypeName: "aws_connect_phone_number",
			Name:     "Phone Number",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceQueue,
			TypeName: "aws_connect_queue",
			Name:     "Queue",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceQuickConnect,
			TypeName: "aws_connect_quick_connect",
			Name:     "Quick Connect",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceRoutingProfile,
			TypeName: "aws_connect_routing_profile",
			Name:     "Routing Profile",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceSecurityProfile,
			TypeName: "aws_connect_security_profile",
			Name:     "Security Profile",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUser,
			TypeName: "aws_connect_user",
			Name:     "User",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUserHierarchyGroup,
			TypeName: "aws_connect_user_hierarchy_group",
			Name:     "User Hierarchy Group",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  resourceUserHierarchyStructure,
			TypeName: "aws_connect_user_hierarchy_structure",
			Name:     "User Hierarchy Structure",
		},
		{
			Factory:  resourceVocabulary,
			TypeName: "aws_connect_vocabulary",
			Name:     "Vocabulary",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Connect
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*connect.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))
	optFns := []func(*connect.Options){
		connect.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
		withExtraOptions(ctx, p, config),
	}

	return connect.NewFromConfig(cfg, optFns...), nil
}

// withExtraOptions returns a functional option that allows this service package to specify extra API client options.
// This option is always called after any generated options.
func withExtraOptions(ctx context.Context, sp conns.ServicePackage, config map[string]any) func(*connect.Options) {
	if v, ok := sp.(interface {
		withExtraOptions(context.Context, map[string]any) []func(*connect.Options)
	}); ok {
		optFns := v.withExtraOptions(ctx, config)

		return func(o *connect.Options) {
			for _, optFn := range optFns {
				optFn(o)
			}
		}
	}

	return func(*connect.Options) {}
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
