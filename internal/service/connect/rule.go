package connect

import (
	"context"
	"fmt"
	"log"
	// "slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/connect"
	awstypes "github.com/aws/aws-sdk-go-v2/service/connect/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	// "github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	// tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	// "github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	// "github.com/hashicorp/terraform-provider-aws/internal/verify"
	// "github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_connect_rule", name="Rule")
func resourceRule() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceRuleCreate,
		ReadWithoutTimeout:   resourceRuleRead,
		UpdateWithoutTimeout: resourceRuleUpdate,
		DeleteWithoutTimeout: resourceRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"event_source_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: enum.Validate[awstypes.EventSourceName](),
			},
			// "integration_association_id": {
			// 	Type:     schema.TypeString,
			// 	Optional: true,
			// },
			// "action_type": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	ValidateDiagFunc: enum.Validate[awstypes.ActionType](),
			// },
			"publish_status": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateDiagFunc: enum.Validate[awstypes.RulePublishStatus](),
			},
			"function": {
				Type:     schema.TypeString,
				Required: true,
			},
			// "trigger_event_source": {
			// 	Type:     schema.TypeSet,
			// 	Required: true,
			// 	MinItems: 0,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"event_source_name": {
			// 				Type:     schema.TypeString,
			// 				Required: true,
			// 				ValidateDiagFunc: enum.Validate[awstypes.EventSourceName](),
			// 			},
			// 			"integration_association_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },
			"actions": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assign_contact_category_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "ASSIGN_CONTACT_CATEGORY",
									},
								},
							},
						},
						"event_bridge_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "GENERATE_EVENTBRIDGE_EVENT",
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"send_notification_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "SEND_NOTIFICATION",
									},
									"delivery_method": {
										Type:     schema.TypeString,
										Required: true,
									},
									"subject": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content": {
										Type:     schema.TypeString,
										Required: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateDiagFunc: enum.Validate[awstypes.NotificationContentType](),
									},
									"recipient": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user_ids": {
													Type:     schema.TypeSet,
													Optional: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"task_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "CREATE_TASK",
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"contact_flow_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"reference": {
										Type:     schema.TypeSet,
										Optional: false,
										Required: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateDiagFunc: enum.Validate[awstypes.ReferenceType](),
												},
											},
										},
									},
								},
							},
						},
						// "update_case_action": {
							// Type:     schema.TypeList,
							// Optional: true,
							// MaxItems: 1,
							// Elem: &schema.Resource{
								// Schema: map[string]*schema.Schema{
									// "fields": {
										// Type: schema.TypeMap,
										// Optional: false
										// Elem: &schema.Schema{
										// Type: schema.TypeString,
										// },
									// },
									// "template_id": {
										// Type:     schema.TypeString,
										// Optional: true,
									// },
								// },
							// },
						// },
						"create_case_action": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_type": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "CREATE_TASK",
									},
									"fields": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)
	var diags diag.Diagnostics
			// IntegrationAssociationId: aws.String(d.Get("trigger_event_source.0.integration_association_id").(string)),

	input := &connect.CreateRuleInput{
		InstanceId: aws.String(d.Get("instance_id").(string)),
		Name:       aws.String(d.Get("name").(string)),
		TriggerEventSource: &awstypes.RuleTriggerEventSource{
			EventSourceName: awstypes.EventSourceName(d.Get("event_source_name").(string)),
		},
		Function: aws.String(d.Get("function").(string)),
		Actions: expandRuleActions(d.Get("actions").([]interface{})),
		PublishStatus: awstypes.RulePublishStatus(d.Get("publish_status").(string)),
	}

			// Actions: []awstypes.RuleAction{}

		// Actions: []awstypes.RuleAction{
		// 	{
		// 		ActionType: awstypes.ActionType(d.Get("action_type").(string)),
		// 	},
		// },

	output, err := conn.CreateRule(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ruleCreateResourceID(d.Get("instance_id").(string), aws.ToString(output.RuleId))

	d.SetId(id)

	return append(diags, resourceRuleRead(ctx, d, meta)...)
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	instanceID, ruleID, err := ruleParseResourceID(d.Id())
	if err != nil {
		return sdkdiag.AppendFromErr(diags, err)
	}

	input := &connect.DescribeRuleInput{
		InstanceId: aws.String(instanceID),
		RuleId:     aws.String(ruleID),
	}

	output, err := conn.DescribeRule(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	if output == nil {
		log.Printf("[WARN] Connect Rule (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", output.Rule.Name)
	d.Set("instance_id", instanceID)
	d.Set("event_source_name", output.Rule.TriggerEventSource.EventSourceName)
	d.Set("function", output.Rule.Function)
	if err := d.Set("actions", flattenActions(output.Rule.Actions)); err != nil {
		return diag.FromErr(err)
	}
	d.Set("publish_status", output.Rule.PublishStatus)

	return diags
}

func flattenActions(actions []awstypes.RuleAction) []interface{} {
	if actions == nil {
		return []interface{}{}
	}

	// tfMap := map[string]interface{}{
	// 	"actions": apiObject.StorageType,
	// }

	// var result []interface{}{}
	tfMap := map[string]interface{}{}

	for _, action := range actions {

		switch action.ActionType {
		case awstypes.ActionTypeAssignContactCategory:
			tfMap["assign_contact_category_action"] = []interface{}{
				map[string]interface{}{
					"action_type": awstypes.ActionTypeAssignContactCategory,
				},
		  }

		case awstypes.ActionTypeGenerateEventbridgeEvent:
			tfMap["event_bridge_action"] = []interface{}{
				map[string]interface{}{
					"action_type": awstypes.ActionTypeGenerateEventbridgeEvent,
					"name": aws.ToString(action.EventBridgeAction.Name),
				},
		  }

		case awstypes.ActionTypeSendNotification:
			tfMap["send_notification_action"] = []interface{}{
				map[string]interface{}{
					"action_type": awstypes.ActionTypeSendNotification,
					"delivery_method": action.SendNotificationAction.DeliveryMethod,
					"subject":         aws.ToString(action.SendNotificationAction.Subject),
					"content":         aws.ToString(action.SendNotificationAction.Content),
					"content_type":    action.SendNotificationAction.ContentType,
					"recipient":       []interface{}{
						map[string]interface{}{
							"user_ids": action.SendNotificationAction.Recipient.UserIds,
						},
					},
				},
			}

			// 	if action.SendNotificationAction.Recipient != nil {
			// 		notificationAction["recipient"] = []interface{}{
			// 			map[string]interface{}{
			// 				"user_ids": action.SendNotificationAction.Recipient.UserIds,
			// 			},
			// 		}
			// 	}
			// }

			// END SWITCH
		}
	}
	return []interface{}{tfMap}
}


// function flattenRecipients(recipient []awstypes.NotificationRecipientType) []interface{} {
// 	recipients := map[string]interface{}{}

// 	recipients["user_ids"] = []interface{}{
// 		map[string]interface{}{
// 			"action_type": awstypes.ActionTypeGenerateEventbridgeEvent,
// 		}

// 	return []interface{}{recipients}
// }

	// return result


	// tfMap["assign_contact_category_action"] = []interface{}{
	// 	map[string]interface{}{
	// 		"action_type": "ASSIGN_CONTACT_CATEGORY",
	// 	},

		// map[string]interface{}{
		// 	"ActionType": "GENERATE_EVENTBRIDGE_EVENT",
		// 	"EventBridgeAction": map[string]interface{}{
		// 		"Name": "tfsample",
		// 	},
		// },



		// case awstypes.ActionTypeGenerateEventbridgeEvent:
		// 	actionMap["action_type"] = awstypes.ActionTypeGenerateEventbridgeEvent
		// 			map[string]interface{}{
		// 				"name": aws.ToString(action.EventBridgeAction.Name),
		// 			}

		// case awstypes.ActionTypeSendNotification:
		// 	actionMap["action_type"] = awstypes.ActionTypeSendNotification
		// 		notificationAction := map[string]interface{}{
		// 			"delivery_method": action.SendNotificationAction.DeliveryMethod,
		// 			"subject":         aws.ToString(action.SendNotificationAction.Subject),
		// 			"content":         aws.ToString(action.SendNotificationAction.Content),
		// 			"content_type":    action.SendNotificationAction.ContentType,
		// 		}

		// 		if action.SendNotificationAction.Recipient != nil {
		// 			notificationAction["recipient"] = []interface{}{
		// 				map[string]interface{}{
		// 					"user_ids": action.SendNotificationAction.Recipient.UserIds,
		// 				},
		// 			}
		// 		}
		// 		actionMap["send_notification_action"] = []interface{}{notificationAction}

	// case awstypes.ActionTypeCreateTask:
	// 	actionMap["action_type"] = awstypes.ActionTypeCreateTask
	// 	if action.TaskAction != nil {
	// 		taskAction := map[string]interface{}{
	// 			"name":           aws.ToString(action.TaskAction.Name),
	// 			"description":    aws.ToString(action.TaskAction.Description),
	// 			"contact_flow_id": aws.ToString(action.TaskAction.ContactFlowId),
	// 			"reference":      flattenTaskReferences(action.TaskAction.References),
	// 		}
	// 	}
	// 	actionMap["task_action"] = []interface{}{taskAction}

	// case awstypes.ActionTypeCreateCase:
	// 	actionMap["action_type"] = awstypes.ActionTypeCreateCase
	// 	if action.CreateCaseAction != nil {
	// 		createCaseAction := map[string]interface{}{
	// 			"fields": flattenCaseFields(action.CreateCaseAction.Fields),
	// 		}

	// 	}


func flattenTaskReferences(references map[string]awstypes.Reference) []interface{} {
	var result []interface{}
	for name, ref := range references {
		result = append(result, map[string]interface{}{
			"name":  name,
			"value": ref.Value,
			"type":  ref.Type,
		})
	}
	return result
}

func flattenCaseFields(fields []awstypes.FieldValue) []interface{} {
	var result []interface{}
	for _, field := range fields {
		result = append(result, map[string]interface{}{
			"id":    field.Id,
			"value": field.Value,
		})
	}
	return result
}

// func flattenActions(actions *awstypes.RuleAction) []interface{} {
// 	if actions == nil {
// 		return []interface{}{}
// 	}

// 	actions := make([]awstypes.RuleAction, 0, len(actions))

// 	tfMap := map[string]interface{}{
// 		names.AttrStorageType: apiObject.StorageType,
// 	}

// }


func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)

	if d.HasChange("name") || d.HasChange("trigger_event_source") || d.HasChange("function") || d.HasChange("actions") || d.HasChange("publish_status") {
		input := &connect.UpdateRuleInput{
			InstanceId: aws.String(d.Get("instance_id").(string)),
			RuleId:     aws.String(d.Id()),
			Name:       aws.String(d.Get("name").(string)),
			// TriggerEventSource: &awstypes.RuleTriggerEventSource{
				// EventSourceName: awstypes.EventSourceName(d.Get("trigger_event_source").([]interface{})[0].(map[string]interface{})["event_source_name"].(string)),
				// IntegrationAssociationId: aws.String(d.Get("trigger_event_source").([]interface{})[0].(map[string]interface{})["integration_association_id"].(string)),
			// },
			Function: aws.String(d.Get("function").(string)),
			Actions: expandRuleActions(d.Get("actions").([]interface{})),
			PublishStatus: awstypes.RulePublishStatus(d.Get("publish_status").(string)),
		}

		_, err := conn.UpdateRule(ctx, input)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return append(diags, resourceRuleRead(ctx, d, meta)...)
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).ConnectClient(ctx)
	var diags diag.Diagnostics

	instanceID, ruleID, err := ruleParseResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	input := &connect.DeleteRuleInput{
		InstanceId: aws.String(instanceID),
		RuleId:    aws.String(ruleID),
	}

	_, error := conn.DeleteRule(ctx, input)
	if error != nil {
		return diag.FromErr(error)
	}

	d.SetId("")

	return diags
}


func expandRuleActions(tfList []interface{}) []awstypes.RuleAction {

	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap, ok := tfList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	actions := make([]awstypes.RuleAction, 0, len(tfList))

	if v, ok := tfMap["assign_contact_category_action"].([]interface{}); ok && len(v) > 0 {
		action := awstypes.RuleAction{
			ActionType: awstypes.ActionTypeAssignContactCategory,
			AssignContactCategoryAction: &awstypes.AssignContactCategoryActionDefinition{},
		}
		actions = append(actions, action)
	}

	if v, ok := tfMap["event_bridge_action"].([]interface{}); ok && len(v) > 0 {
		action := awstypes.RuleAction{
			ActionType: awstypes.ActionTypeGenerateEventbridgeEvent,
			EventBridgeAction: &awstypes.EventBridgeActionDefinition{
				Name: aws.String(v[0].(map[string]interface{})["name"].(string)),
			},
		}
		actions = append(actions, action)
	}

	if v, ok := tfMap["send_notification_action"].([]interface{}); ok && len(v) > 0 {
		action := awstypes.RuleAction{
			ActionType: awstypes.ActionTypeSendNotification,
			SendNotificationAction: &awstypes.SendNotificationActionDefinition{
				DeliveryMethod: awstypes.NotificationDeliveryType(v[0].(map[string]interface{})["delivery_method"].(string)),
				Subject:        aws.String(v[0].(map[string]interface{})["subject"].(string)),
				Content:        aws.String(v[0].(map[string]interface{})["content"].(string)),
				ContentType:    awstypes.NotificationContentType(v[0].(map[string]interface{})["content_type"].(string)),
				Recipient: &awstypes.NotificationRecipientType{
					UserIds: expandStringSet(v[0].(map[string]interface{})["recipient"].([]interface{})[0].(map[string]interface{})["user_ids"].(*schema.Set)),
			}},
		}
		actions = append(actions, action)
	}

	return actions
}

func expandStringSet(s *schema.Set) []string {
	if s == nil {
		return nil
	}
	result := make([]string, 0, s.Len())
	for _, v := range s.List() {
		result = append(result, v.(string))
	}
	return result
}


			// if recipient, ok := notificationAction["recipient"].([]interface{}); ok && len(recipient) > 0 {
				// recipientData := recipient[0].(map[string]interface{})
				// action.SendNotificationAction.Recipient = &awstypes.NotificationRecipientType{
					// UserIds: aws.String(recipientData.(map[string]interface{})["user_ids"].(string)),
				// }
			// }
// 		actions = append(actions, action)
// 	}

// 	return actions
// }


















// func expandRuleActions(tfList []interface{}) []awstypes.RuleAction {

// 	if len(tfList) == 0 || tfList[0] == nil {
// 		return nil
// 	}

// 	actions := make([]awstypes.RuleAction, 0, len(tfList))

// 	for _, item := range tfList {
// 		data, ok := item.(map[string]interface{})
// 		if !ok {
// 			continue
// 		}
// 		action := awstypes.RuleAction{
// 			ActionType: awstypes.ActionType(data["action_type"].(string)),
// 			// ActionType: awstypes.ActionTypeAssignContactCategory,
// 		}

// 		switch action.ActionType {
// 		case awstypes.ActionTypeAssignContactCategory:
// 			// action.ActionType = awstypes.ActionTypeAssignContactCategory
// 			action.AssignContactCategoryAction = &awstypes.AssignContactCategoryActionDefinition{}

// 		case awstypes.ActionTypeGenerateEventbridgeEvent:
// 			if v, ok := data["event_bridge_action"].([]interface{}); ok && len(v) > 0 {
// 				// action.ActionType = awstypes.ActionTypeGenerateEventbridgeEvent
// 				action.EventBridgeAction = &awstypes.EventBridgeActionDefinition{
// 					Name: aws.String(v[0].(map[string]interface{})["name"].(string)),
// 				}
// 			}

// 		}

// 		actions = append(actions, action)
// 	}

// 	return actions
// }



		// case awstypes.ActionTypeSendNotification:
		// 	if v, ok := data["send_notification_action"].([]interface{}); ok && len(v) > 0 {
		// 		notificationAction := v[0].(map[string]interface{})
		// 		action.SendNotificationAction = &awstypes.SendNotificationActionDefinition{
		// 			DeliveryMethod: awstypes.NotificationDeliveryType(notificationAction["delivery_method"].(string)),
		// 			Subject:        aws.String(notificationAction["subject"].(string)),
		// 			Content:        aws.String(notificationAction["content"].(string)),
		// 			ContentType:    awstypes.NotificationContentType(notificationAction["content_type"].(string)),
		// 		}
		// 		if recipient, ok := notificationAction["recipient"].([]interface{}); ok && len(recipient) > 0 {
		// 			recipientData := recipient[0].(map[string]interface{})
		// 			action.SendNotificationAction.Recipient = &awstypes.NotificationRecipientType{
		// 				UserIds:        expandStringList(recipientData["user_arns"].([]interface{})),
		// 				// QueueArns:       expandStringList(recipientData["queue_arns"].([]interface{})),
		// 				// ContactFlowArns: expandStringList(recipientData["contact_flow_arns"].([]interface{})),
		// 				// EmailAddresses:  expandStringList(recipientData["email_addresses"].([]interface{})),
		// 				// SmsPhoneNumbers: expandStringList(recipientData["sms_phone_numbers"].([]interface{})),
		// 			}
		// 		}
		// 	}
		// case awstypes.TaskAction:
		// 	if v, ok := data["task_action"].([]interface{}); ok && len(v) > 0 {
		// 		taskAction := v[0].(map[string]interface{})
		// 		action.TaskAction = awstypes.TaskActionDefinition{
		// 			Name:           aws.String(taskAction["name"].(string)),
		// 			Description:    aws.String(taskAction["description"].(string)),
		// 			ContactFlowArn: aws.String(taskAction["contact_flow_arn"].(string)),
		// 			References:     expandRuleActionParameters(taskAction["references"].(map[string]interface{})),
		// 		}
		// 	}
		// case awstypes.ActionTypeUpdateCase:
		// 	if v, ok := data["update_case_action"].([]interface{}); ok && len(v) > 0 {
		// 		updateCaseAction := v[0].(map[string]interface{})
		// 		action.UpdateCaseAction = &awstypes.UpdateCaseActionDefinition{
		// 			CaseId: aws.String(updateCaseAction["case_id"].(string)),
		// 			Fields: expandRuleActionParameters(updateCaseAction["fields"].(map[string]interface{})),
		// 		}
		// 	}
		// }


// func expandRuleActionParameters(data map[string]interface{}) map[string]string {
// 	parameters := make(map[string]string)

// 	for k, v := range data {
// 		if k != "action_type" {
// 			parameters[k] = v.(string)
// 		}
// 	}

// 	return parameters
// }

// func flattenRuleActions(actions []awstypes.RuleAction) []interface{} {
// 	flattened := make([]interface{}, 0, len(actions))

// 	for _, action := range actions {
// 		data := map[string]interface{}{
// 			"action_type": action.ActionType,
// 		}
// 		// for k, v := range action.Parameters {
// 		// 	data[k] = v
// 		// }
// 		flattened = append(flattened, data)
// 	}

// 	return flattened
// }

// func flattenRuleTriggerEventSource(triggerEventSource *awstypes.RuleTriggerEventSource) []interface{} {
// 	if triggerEventSource == nil {
// 		return nil
// 	}

// 	return []interface{}{
// 		map[string]interface{}{
// 			"event_source_name":          triggerEventSource.EventSourceName,
// 			"integration_association_id": aws.ToString(triggerEventSource.IntegrationAssociationId),
// 		},
// 	}
// }

func ruleCreateResourceID(instanceID, ruleID string) string {
	parts := []string{instanceID, ruleID}
	id := strings.Join(parts, ",")

	return id
}

func ruleParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, ",", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected instanceID,ruleID", id)
	}

	return parts[0], parts[1], nil
}