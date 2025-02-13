package connect

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	awstypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	phoneNumberContactFlowResourceIDSeparator = ","
)
// @SDKResource("aws_connect_associate_phone_number_contact_flow", name="Associate Phone Number Contact Flow")
func resourceAssociatePhoneNumberContactFlow() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAssociatePhoneNumberContactFlowCreate,
		ReadWithoutTimeout:   resourceAssociatePhoneNumberContactFlowRead,
		DeleteWithoutTimeout: resourceAssociatePhoneNumberContactFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"phone_number_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contact_flow_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAssociatePhoneNumberContactFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN := d.Get("instance_arn").(string)
	phoneNumberARN := d.Get("phone_number_arn").(string)
	contactFlowARN := d.Get("contact_flow_arn").(string)

	input := &connect.AssociatePhoneNumberContactFlowInput{
		InstanceId:     aws.String(instanceARN),
		PhoneNumberId:  aws.String(phoneNumberARN),
		ContactFlowId:  aws.String(contactFlowARN),
	}

	_, err := conn.AssociatePhoneNumberContactFlow(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	id := phoneNumberContactFlowCreateResourceID(instanceARN, phoneNumberARN)
	d.SetId(id)

	return append(diags, resourceAssociatePhoneNumberContactFlowRead(ctx, d, meta)...)
}

func resourceAssociatePhoneNumberContactFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN, phoneNumberARN, err := phoneNumberContactFlowParseResourceID(d.Id())

  contactFlowARN, err := findAssociatedContactFlow(ctx, conn, instanceARN, phoneNumberARN)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Connect Instance Storage Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	d.Set("instance_arn", instanceARN)
	d.Set("phone_number_arn", phoneNumberARN)
	d.Set("contact_flow_arn", contactFlowARN)

	return diags
}

func findAssociatedContactFlow(ctx context.Context, conn *connect.Client, instanceARN string, phoneNumberARN string) (string, error) {
	input := &connect.ListFlowAssociationsInput{
		InstanceId: aws.String(instanceARN),
		ResourceType: awstypes.ListFlowAssociationResourceTypeVoicePhoneNumber,
	}

	output, err := conn.ListFlowAssociations(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return "", &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return "", err
	}

	if output == nil || output.FlowAssociationSummaryList == nil {
		return "", tfresource.NewEmptyResultError(input)
	}

  var flowAssociation string
	for _, v := range output.FlowAssociationSummaryList {
		if *v.ResourceId == phoneNumberARN {
			flowAssociation = *v.FlowId
		}
	}

	return flowAssociation, nil
}

func resourceAssociatePhoneNumberContactFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN := d.Get("instance_arn").(string)
	phoneNumberARN := d.Get("phone_number_arn").(string)

	input := &connect.DisassociatePhoneNumberContactFlowInput{
		InstanceId:    aws.String(instanceARN),
		PhoneNumberId: aws.String(phoneNumberARN),
	}

	_, err := conn.DisassociatePhoneNumberContactFlow(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func phoneNumberContactFlowCreateResourceID(instanceARN, phoneNumberARN string) string {
	parts := []string{instanceARN, phoneNumberARN}
	id := strings.Join(parts, phoneNumberContactFlowResourceIDSeparator)

	return id
}

func phoneNumberContactFlowParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, phoneNumberContactFlowResourceIDSeparator, 3)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%[1]s), expected instanceARN%[2]sphoneNumberARN", id, phoneNumberContactFlowResourceIDSeparator)
	}

	return parts[0], parts[1], nil
}