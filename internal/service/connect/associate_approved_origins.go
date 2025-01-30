package connect

import (
	"context"
	"fmt"
	// "log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	// awstypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	// "github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	associateApprovedOriginIDSeparator = ","
)
// @SDKResource("aws_connect_associate_approved_origin", name="Associate Approved Origin")
func resourceAssociateApprovedOrigin() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAssociateApprovedOriginCreate,
		ReadWithoutTimeout:   resourceAssociateApprovedOriginRead,
		DeleteWithoutTimeout: resourceAssociateApprovedOriginDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_arn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"origin": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAssociateApprovedOriginCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN := d.Get("instance_arn").(string)
	origin := d.Get("origin").(string)

	input := &connect.AssociateApprovedOriginInput{
		InstanceId: aws.String(instanceARN),
		Origin:     aws.String(origin),
	}

	_, err := conn.AssociateApprovedOrigin(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	id := approvedOriginCreateResourceID(instanceARN, origin)
	d.SetId(id)

	return append(diags, resourceAssociateApprovedOriginRead(ctx, d, meta)...)
}

func resourceAssociateApprovedOriginRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN, origin, err := approvedOriginParseResourceID(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	origins, err := findAssociatedApprovedOrigin(ctx, conn, instanceARN)

	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for _, o := range origins {
		if o == origin {
			found = true
			break
		}
	}

	if !found {
		d.SetId("")
		return diags
	}

	d.Set("instance_arn", instanceARN)
	d.Set("origin", origin)

	return diags
}

func findAssociatedApprovedOrigin(ctx context.Context, conn *connect.Client, instanceARN string) ([]string, error) {
	input := &connect.ListApprovedOriginsInput{
		InstanceId: aws.String(instanceARN),
	}

	output, err := conn.ListApprovedOrigins(ctx, input)

	if err != nil {
		return nil, err
	}

	if output == nil || output.Origins == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	var origins []string
	for _, v := range output.Origins {
		origins = append(origins, v)
	}

	return origins, nil
}

func resourceAssociateApprovedOriginDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceARN := d.Get("instance_arn").(string)
	origin := d.Get("origin").(string)

	input := &connect.DisassociateApprovedOriginInput{
		InstanceId: aws.String(instanceARN),
		Origin:     aws.String(origin),
	}

	_, err := conn.DisassociateApprovedOrigin(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func approvedOriginCreateResourceID(instanceARN, origin string) string {
	parts := []string{instanceARN, origin}
	id := strings.Join(parts, associateApprovedOriginIDSeparator)

	return id
}

func approvedOriginParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, associateApprovedOriginIDSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%[1]s), expected instanceARN%[2]sorigin", id, associateApprovedOriginIDSeparator)
	}

	return parts[0], parts[1], nil
}