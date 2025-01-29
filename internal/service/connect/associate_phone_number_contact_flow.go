package connect

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	// awstypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	// "github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	phoneNumberContactFlowResourceIDSeparator = ":"
)
// @SDKResource("aws_connect_associate_phone_number_contact_flow", name="Associate Phone Number Contact Flow")
func resourceAssociatePhoneNumberContactFlow() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAssociatePhoneNumberContactFlowCreate,
		ReadWithoutTimeout:   resourceAssociatePhoneNumberContactFlowRead,
		DeleteWithoutTimeout: resourceAssociatePhoneNumberContactFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePhoneNumberContactFlowImport,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"phone_number_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contact_flow_id": {
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

	instanceID := d.Get("instance_id").(string)
	phoneNumberID := d.Get("phone_number_id").(string)
	contactFlowID := d.Get("contact_flow_id").(string)

	input := &connect.AssociatePhoneNumberContactFlowInput{
		InstanceId:     aws.String(instanceID),
		PhoneNumberId:  aws.String(phoneNumberID),
		ContactFlowId:  aws.String(contactFlowID),
	}

	_, err := conn.AssociatePhoneNumberContactFlow(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	id := phoneNumberContactFlowCreateResourceID(instanceID, phoneNumberID, contactFlowID)
	d.SetId(id)

	return append(diags, resourceAssociatePhoneNumberContactFlowRead(ctx, d, meta)...)
}

func resourceAssociatePhoneNumberContactFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// instanceID := d.Get("instance_id").(string)
	// phoneNumberID := d.Get("phone_number_id").(string)
	// contact_flow_id := d.Get("contact_flow_id").(string)

	instanceID, phoneNumberID, contactFlowID, err := phoneNumberContactFlowParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Connect Instance Storage Config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	d.Set("instance_id", instanceID)
	d.Set("phone_number_id", phoneNumberID)
	d.Set("contact_flow_id", contactFlowID)

	return diags
}

func resourceAssociatePhoneNumberContactFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceID := d.Get("instance_id").(string)
	phoneNumberID := d.Get("phone_number_id").(string)

	input := &connect.DisassociatePhoneNumberContactFlowInput{
		InstanceId:    aws.String(instanceID),
		PhoneNumberId: aws.String(phoneNumberID),
	}

	_, err := conn.DisassociatePhoneNumberContactFlow(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func resourcePhoneNumberContactFlowImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), phoneNumberContactFlowResourceIDSeparator, 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%[1]s), expected instanceID%[2]sphoneNumberID%[2]scontactFlowID", parts, phoneNumberContactFlowResourceIDSeparator)
	}
	instanceID := parts[0]
	phoneNumberID := parts[1]
	contactFlowID := parts[2]
	d.Set("instance_id", instanceID)
	d.Set("phone_number_id", phoneNumberID)
	d.Set("contact_flow_id", contactFlowID)
	d.SetId(fmt.Sprintf("%s:%s:%s", instanceID, phoneNumberID, contactFlowID))
	return []*schema.ResourceData{d}, nil
}

func phoneNumberContactFlowCreateResourceID(instanceID, phoneNumberID string, contactFlowID string) string {
	parts := []string{instanceID, phoneNumberID, contactFlowID}
	id := strings.Join(parts, phoneNumberContactFlowResourceIDSeparator)

	return id
}

func phoneNumberContactFlowParseResourceID(id string) (string, string, string, error) {
	parts := strings.SplitN(id, phoneNumberContactFlowResourceIDSeparator, 3)
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", "", "", fmt.Errorf("unexpected format of ID (%[1]s), expected instanceID%[2]sphoneNumberID%[2]scontactFlowID", id, phoneNumberContactFlowResourceIDSeparator)
	}

	return parts[0], parts[1], parts[2], nil
}