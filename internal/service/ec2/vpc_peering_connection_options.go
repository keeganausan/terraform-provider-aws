package ec2

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func ResourceVPCPeeringConnectionOptions() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCPeeringConnectionOptionsCreate,
		Read:   resourceVPCPeeringConnectionOptionsRead,
		Update: resourceVPCPeeringConnectionOptionsUpdate,
		Delete: resourceVPCPeeringConnectionOptionsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"accepter":  vpcPeeringConnectionOptionsSchema,
			"requester": vpcPeeringConnectionOptionsSchema,
			"vpc_peering_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVPCPeeringConnectionOptionsCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("vpc_peering_connection_id").(string))

	return resourceVPCPeeringConnectionOptionsUpdate(d, meta)
}

func resourceVPCPeeringConnectionOptionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	pc, err := vpcPeeringConnection(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error reading VPC Peering Connection (%s): %w", d.Id(), err)
	}

	if pc == nil {
		log.Printf("[WARN] VPC Peering Connection (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("vpc_peering_connection_id", pc.VpcPeeringConnectionId)

	if err := d.Set("accepter", flattenVPCPeeringConnectionOptions(pc.AccepterVpcInfo.PeeringOptions)); err != nil {
		return fmt.Errorf("error setting VPC Peering Connection Options accepter information: %s", err)
	}
	if err := d.Set("requester", flattenVPCPeeringConnectionOptions(pc.RequesterVpcInfo.PeeringOptions)); err != nil {
		return fmt.Errorf("error setting VPC Peering Connection Options requester information: %s", err)
	}

	return nil
}

func resourceVPCPeeringConnectionOptionsUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*conns.AWSClient).EC2Conn

	pc, err := vpcPeeringConnection(conn, d.Id())

	if err != nil {
		return fmt.Errorf("error reading VPC Peering Connection (%s): %w", d.Id(), err)
	}

	if pc == nil {
		return fmt.Errorf("VPC Peering Connection (%s) not found", d.Id())
	}

	if d.HasChanges("accepter", "requester") {
		crossRegionPeering := aws.StringValue(pc.RequesterVpcInfo.Region) != aws.StringValue(pc.AccepterVpcInfo.Region)

		input := &ec2.ModifyVpcPeeringConnectionOptionsInput{
			VpcPeeringConnectionId: aws.String(d.Id()),
		}
		if d.HasChange("accepter") {
			input.AccepterPeeringConnectionOptions = expandVPCPeeringConnectionOptions(d.Get("accepter").([]interface{}), crossRegionPeering)
		}
		if d.HasChange("requester") {
			input.RequesterPeeringConnectionOptions = expandVPCPeeringConnectionOptions(d.Get("requester").([]interface{}), crossRegionPeering)
		}

		log.Printf("[DEBUG] Modifying VPC Peering Connection options: %s", input)
		_, err = conn.ModifyVpcPeeringConnectionOptions(input)

		if err != nil {
			return fmt.Errorf("error modifying VPC Peering Connection (%s) Options: %w", d.Id(), err)
		}

		// Retry reading back the modified options to deal with eventual consistency.
		// Often this is to do with a delay transitioning from pending-acceptance to active.
		err = resource.Retry(3*time.Minute, func() *resource.RetryError {
			pc, err = vpcPeeringConnection(conn, d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}

			if pc == nil {
				return nil
			}

			if d.HasChange("accepter") && pc.AccepterVpcInfo != nil {
				if aws.BoolValue(pc.AccepterVpcInfo.PeeringOptions.AllowDnsResolutionFromRemoteVpc) != aws.BoolValue(input.AccepterPeeringConnectionOptions.AllowDnsResolutionFromRemoteVpc) ||
					aws.BoolValue(pc.AccepterVpcInfo.PeeringOptions.AllowEgressFromLocalClassicLinkToRemoteVpc) != aws.BoolValue(input.AccepterPeeringConnectionOptions.AllowEgressFromLocalClassicLinkToRemoteVpc) ||
					aws.BoolValue(pc.AccepterVpcInfo.PeeringOptions.AllowEgressFromLocalVpcToRemoteClassicLink) != aws.BoolValue(input.AccepterPeeringConnectionOptions.AllowEgressFromLocalVpcToRemoteClassicLink) {
					return resource.RetryableError(fmt.Errorf("VPC Peering Connection (%s) accepter Options not stable", d.Id()))
				}
			}
			if d.HasChange("requester") && pc.RequesterVpcInfo != nil {
				if aws.BoolValue(pc.RequesterVpcInfo.PeeringOptions.AllowDnsResolutionFromRemoteVpc) != aws.BoolValue(input.RequesterPeeringConnectionOptions.AllowDnsResolutionFromRemoteVpc) ||
					aws.BoolValue(pc.RequesterVpcInfo.PeeringOptions.AllowEgressFromLocalClassicLinkToRemoteVpc) != aws.BoolValue(input.RequesterPeeringConnectionOptions.AllowEgressFromLocalClassicLinkToRemoteVpc) ||
					aws.BoolValue(pc.RequesterVpcInfo.PeeringOptions.AllowEgressFromLocalVpcToRemoteClassicLink) != aws.BoolValue(input.RequesterPeeringConnectionOptions.AllowEgressFromLocalVpcToRemoteClassicLink) {
					return resource.RetryableError(fmt.Errorf("VPC Peering Connection (%s) requester Options not stable", d.Id()))
				}
			}

			return nil
		})
	}

	return resourceVPCPeeringConnectionOptionsRead(d, meta)
}

func resourceVPCPeeringConnectionOptionsDelete(d *schema.ResourceData, meta interface{}) error {
	// Don't do anything with the underlying VPC peering connection.
	return nil
}

// vpcPeeringConnection returns the VPC peering connection corresponding to the specified identifier.
// Returns nil if no VPC peering connection is found or the connection has reached a terminal state
// according to https://docs.aws.amazon.com/vpc/latest/peering/vpc-peering-basics.html#vpc-peering-lifecycle.
func vpcPeeringConnection(conn *ec2.EC2, vpcPeeringConnectionID string) (*ec2.VpcPeeringConnection, error) {
	outputRaw, _, err := StatusVPCPeeringConnectionDeleted(conn, vpcPeeringConnectionID)()

	if output, ok := outputRaw.(*ec2.VpcPeeringConnection); ok {
		return output, err
	}

	return nil, err
}

func expandVPCPeeringConnectionOptions(vOptions []interface{}, crossRegionPeering bool) *ec2.PeeringConnectionOptionsRequest {
	if len(vOptions) == 0 || vOptions[0] == nil {
		return nil
	}

	mOptions := vOptions[0].(map[string]interface{})

	options := &ec2.PeeringConnectionOptionsRequest{}

	if v, ok := mOptions["allow_remote_vpc_dns_resolution"].(bool); ok {
		options.AllowDnsResolutionFromRemoteVpc = aws.Bool(v)
	}
	if !crossRegionPeering {
		if v, ok := mOptions["allow_classic_link_to_remote_vpc"].(bool); ok {
			options.AllowEgressFromLocalClassicLinkToRemoteVpc = aws.Bool(v)
		}
		if v, ok := mOptions["allow_vpc_to_remote_classic_link"].(bool); ok {
			options.AllowEgressFromLocalVpcToRemoteClassicLink = aws.Bool(v)
		}
	}

	return options
}

func flattenVPCPeeringConnectionOptions(options *ec2.VpcPeeringConnectionOptionsDescription) []interface{} {
	// When the VPC Peering Connection is pending acceptance,
	// the details about accepter and/or requester peering
	// options would not be included in the response.
	if options == nil {
		return []interface{}{}
	}

	return []interface{}{map[string]interface{}{
		"allow_remote_vpc_dns_resolution":  aws.BoolValue(options.AllowDnsResolutionFromRemoteVpc),
		"allow_classic_link_to_remote_vpc": aws.BoolValue(options.AllowEgressFromLocalClassicLinkToRemoteVpc),
		"allow_vpc_to_remote_classic_link": aws.BoolValue(options.AllowEgressFromLocalVpcToRemoteClassicLink),
	}}
}
