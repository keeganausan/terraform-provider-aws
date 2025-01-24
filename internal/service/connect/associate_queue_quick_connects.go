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
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	// "github.com/hashicorp/terraform-provider-aws/internal/flex"
	// tfslices "github.com/hashicorp/terraform-provider-aws/internal/slices"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"

)

const (
	maxResults = 100
	queueQuickConnectsResourceIDPartCount = 2
	queueQuickConnectsResourceIDSeparator = ","
)

// @SDKResource("aws_connect_associate_queue_quick_connects", name="Associate Queue Quick Connects")
func resourceAssociateQueueQuickConnects() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceAssociateQueueQuickConnectsCreate,
		ReadWithoutTimeout:   resourceAssociateQueueQuickConnectsRead,
		DeleteWithoutTimeout: resourceAssociateQueueQuickConnectsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quick_connect_ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAssociateQueueQuickConnectsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceID := d.Get("instance_id").(string)
	queueID := d.Get("queue_id").(string)
	quickConnectIDs := expandStringList(d.Get("quick_connect_ids").([]interface{}))

	input := &connect.AssociateQueueQuickConnectsInput{
		InstanceId:      aws.String(instanceID),
		QueueId:         aws.String(queueID),
		QuickConnectIds: quickConnectIDs,
	}
  id := queueQuickConnectsCreateResourceID(instanceID, queueID)
	d.SetId(id)
  log.Printf("[WARN] Connect Queue Quick Connect Associations id (%s) not found, given (%s) ", d.Id(), id)
	log.Printf("[DEBUG] Associating Queue Quick Connect: %s", d.Id())
	_, err := conn.AssociateQueueQuickConnects(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}


	return append(diags, resourceAssociateQueueQuickConnectsRead(ctx, d, meta)...)
}


func resourceAssociateQueueQuickConnectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)
	var queueQuickConnectIDs []string

	instanceID, queueID, err := queueQuickConnectsParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	queueQuickConnects, err := findQueueQuickConnectsByTwoPartKey(ctx, conn, instanceID, queueID)

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Connect Queue Quick Connect Associations (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Connect Queue Quick Connects (%s): %s", d.Id(), err)
	}

	if len(queueQuickConnects) == 0 {
		return sdkdiag.AppendErrorf(diags, "Found Connect Queue Quick Connect Configuration for queue (%s) but there are no QuickConnectIds associated with it. QuickConnectIds must include at least 1 item", queueID)
	}

	d.Set("instance_id", instanceID)
	d.Set("queue_id", queueID)

	for _, v := range queueQuickConnects {
		queueQuickConnectIDs = append(queueQuickConnectIDs, aws.ToString(v.Id))
	}
	d.Set("quick_connect_ids", queueQuickConnectIDs)

	return diags
}

func findQueueQuickConnectsByTwoPartKey(ctx context.Context, conn *connect.Client, instanceID string, queueID string) ([]awstypes.QuickConnectSummary, error) {

	input := &connect.ListQueueQuickConnectsInput{
		InstanceId: aws.String(instanceID),
		QueueId:    aws.String(queueID),
		MaxResults: aws.Int32(maxResults),
	}

	return findQueueQuickConnects(ctx, conn, input)

}

func queueQuickConnectsCreateResourceID(instanceID, queueID string) string {
	parts := []string{instanceID, queueID}
	id := strings.Join(parts, queueQuickConnectsResourceIDSeparator)

	return id
}

func findQueueQuickConnects(ctx context.Context, conn *connect.Client, input *connect.ListQueueQuickConnectsInput) ([]awstypes.QuickConnectSummary, error) {
	output, err := conn.ListQueueQuickConnects(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.QuickConnectSummaryList == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.QuickConnectSummaryList, nil
}











	// parts, err := flex.ExpandResourceId(d.Id(), queueQuickConnectsResourceIDPartCount, true)
	// if err != nil {
	// 	return sdkdiag.AppendFromErr(diags, err)
	// }

	// instanceID, queueID := parts[0], parts[1]
	// _, err = findQueueQuickConnectAssociationByKeys(ctx, conn, instanceID, queueID, func(v *awstypes.QuickConnectSummary))

	// if !d.IsNewResource() && tfresource.NotFound(err) {
	// 	log.Printf("[WARN] Connect Queue Quick Connect Association (%s) not found, removing from state", d.Id())
	// 	d.SetId("")
	// 	return diags
	// }

	// if err != nil {
	// 	return sdkdiag.AppendErrorf(diags, "reading Connect Queue Quick Connects (%s): %s", d.Id(), err)
	// }

	// d.Set("instance_id", instanceID)
	// d.Set("queue_id", queueID)


	// instanceID := d.Get("instance_id").(string)
	// queueID := d.Get("queue_id").(string)

	// input := &connect.DescribeQueueInput{
	// 	InstanceId: aws.String(instanceID),
	// 	QueueId:    aws.String(queueID),
	// }

	// _, err := conn.DescribeQueue(ctx, input)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }



// func findQueueQuickConnectAssociationByKeys(ctx context.Context, conn *connect.Client, instanceID string, queueID string, filter tfslices.Predicate[*awstypes.QuickConnectSummary]) (*string, error) {
// 	var output []string
// 	const maxResults = 100

// 	input := &connect.ListQueueQuickConnectsInput{
// 		InstanceId: aws.String(instanceID),
// 		QueueId:    aws.String(queueID),
// 		MaxResults: aws.Int32(maxResults),
// 	}

// 	pages := connect.NewListQueueQuickConnectsPaginator(conn, input)

// 	for pages.HasMorePages() {
// 		page, err := pages.NextPage(ctx)

// 		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
// 			return nil, &retry.NotFoundError{
// 				LastError:   err,
// 				LastRequest: input,
// 			}
// 		}

// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, v := range page.QuickConnectSummaryList {
// 			if filter(v) {
// 				output = append(output, v)
// 			}
// 		}
// 	}

// 	return output, nil
// }


	// return findQueueQuickConnectAssociation(ctx, conn, input, func(v *awstypes.QuickConnectSummary) bool {
	// 	return v == QuickConnectSummaryList
	// })
// }


// 	return findQueueQuickConnect(ctx, conn, input, func(v *awstypes.QuickConnectConfig) bool {
// 		// return v.QuickConnect. == queueID
// 		return aws.ToString(v.QueueId) == queueID
// 	})
// }

// func findQueueQuickConnect(ctx context.Context, conn *connect.Client, input *connect.ListQueueQuickConnectsInput, filter tfslices.Predicate[*awstypes.QueueQuickConnectConfig]) (*awstypes.QuickConnect, error) {
// // func findQueueQuickConnect(ctx context.Context, conn *connect.Client, input *connect.ListQueueQuickConnectsInput, filter tfslices.Predicate[*awstypes.QueueQuickConnectConfig]) (*awstypes.QuickConnect, error) {
// 	output, err := findQueueQuickConnects(ctx, conn, input, filter)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return tfresource.AssertSingleValueResult(output)
// }

// func findQueueQuickConnects(ctx context.Context, conn *connect.Client, input *connect.ListQueueQuickConnectsInput, filter tfslices.Predicate[*awstypes.QuickConnectConfig]) (*awstypes.QuickConnect, error) {
// 	var output []string
// 	pages := connect.NewListQueueQuickConnectsPaginator(conn, input)
// 	for pages.HasMorePages() {
// 		page, err := pages.NextPage(ctx)

// 		if errs.IsA[*awstypes.ResourceNotFoundException](err) {
// 			return nil, &retry.NotFoundError{
// 				LastError:   err,
// 				LastRequest: input,
// 			}
// 		}

// 		if err != nil {
// 			return nil, err
// 		}

// 		for _, v := range page.QuickConnectSummaryList {
// 			if filter(&v) {
// 				output = append(output, v)
// 			}
// 		}
// 	}

// 	return output, nil
// }

func resourceAssociateQueueQuickConnectsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceID := d.Get("instance_id").(string)
	queueID := d.Get("queue_id").(string)
	quickConnectIDs := expandStringList(d.Get("quick_connect_ids").([]interface{}))

	input := &connect.DisassociateQueueQuickConnectsInput{
		InstanceId: aws.String(instanceID),
		QueueId:    aws.String(queueID),
		QuickConnectIds: quickConnectIDs,
	}

	_, err := conn.DisassociateQueueQuickConnects(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func expandStringList(list []interface{}) []string {
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}

func queueQuickConnectsResourceID(instanceID, queueID string) string {
	parts := []string{instanceID, queueID}
	id := strings.Join(parts, queueQuickConnectsResourceIDSeparator)

	return id
}

func queueQuickConnectsParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, queueQuickConnectsResourceIDSeparator, 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%[1]s), expected instanceID%[2]squeueID", id, queueQuickConnectsResourceIDSeparator)
	}

	return parts[0], parts[1], nil
}