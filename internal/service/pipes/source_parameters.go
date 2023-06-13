package pipes

import (
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/pipes/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func sourceParametersSchema() *schema.Schema {
	verifySecretsManagerARN := validation.StringMatch(regexp.MustCompile(`^(^arn:aws([a-z]|\-)*:secretsmanager:([a-z]{2}((-gov)|(-iso(b?)))?-[a-z]+-\d{1}):(\d{12}):secret:.+)$`), "")

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"activemq_broker_parameters": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.self_managed_kafka",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"credentials": {
								Type:     schema.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"basic_auth": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"queue_name": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 1000),
									validation.StringMatch(regexp.MustCompile(`^[\s\S]*$`), ""),
								),
							},
						},
					},
				},
				"dynamodb_stream_parameters": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.self_managed_kafka",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"dead_letter_config": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"arn": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verify.ValidARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"maximum_record_age_in_seconds": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.Any(
									validation.IntInSlice([]int{-1}),
									validation.IntBetween(60, 604_800),
								),
							},
							"maximum_retry_attempts": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(-1, 10_000),
							},
							"on_partial_batch_item_failure": {
								Type:             schema.TypeString,
								Optional:         true,
								ValidateDiagFunc: enum.Validate[types.OnPartialBatchItemFailureStreams](),
							},
							"parallelization_factor": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10),
								Default:      1,
							},
							"starting_position": {
								Type:             schema.TypeString,
								Required:         true,
								ForceNew:         true,
								ValidateDiagFunc: enum.Validate[types.DynamoDBStreamStartPosition](),
							},
						},
					},
				},
				"filter_criteria": {
					Type:             schema.TypeList,
					Optional:         true,
					MaxItems:         1,
					DiffSuppressFunc: verify.SuppressMissingOptionalConfigurationBlock,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"filter": {
								Type:     schema.TypeList,
								Required: true,
								MaxItems: 5,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"pattern": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringLenBetween(1, 4096),
										},
									},
								},
							},
						},
					},
				},
				"kinesis_stream_parameters": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.self_managed_kafka",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"dead_letter_config": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"arn": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verify.ValidARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"maximum_record_age_in_seconds": {
								Type:     schema.TypeInt,
								Optional: true,
								Computed: true,
								ValidateFunc: validation.Any(
									validation.IntInSlice([]int{-1}),
									validation.IntBetween(60, 604_800),
								),
							},
							"maximum_retry_attempts": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(-1, 10_000),
							},
							"on_partial_batch_item_failure": {
								Type:             schema.TypeString,
								Optional:         true,
								ValidateDiagFunc: enum.Validate[types.OnPartialBatchItemFailureStreams](),
							},
							"parallelization_factor": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10),
								Default:      1,
							},
							"starting_position": {
								Type:             schema.TypeString,
								Required:         true,
								ForceNew:         true,
								ValidateDiagFunc: enum.Validate[types.KinesisStreamStartPosition](),
							},
							"starting_position_timestamp": {
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.IsRFC3339Time,
							},
						},
					},
				},
				"managed_streaming_kafka_parameters": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.self_managed_kafka",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"consumer_group_id": {
								Type:     schema.TypeString,
								Optional: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 200),
									validation.StringMatch(regexp.MustCompile(`^[^.]([a-zA-Z0-9\-_.]+)$`), ""),
								),
							},
							"credentials": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"client_certificate_tls_auth": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
										"sasl_scram_512_auth": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"starting_position": {
								Type:             schema.TypeString,
								Optional:         true,
								ForceNew:         true,
								ValidateDiagFunc: enum.Validate[types.MSKStartPosition](),
							},
							"topic_name": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 249),
									validation.StringMatch(regexp.MustCompile(`^[^.]([a-zA-Z0-9\-_.]+)$`), ""),
								),
							},
						},
					},
				},
				"rabbit_mq_broker": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.self_managed_kafka",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"credentials": {
								Type:     schema.TypeList,
								Required: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"basic_auth": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"queue": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 1000),
									validation.StringMatch(regexp.MustCompile(`^[\s\S]*$`), ""),
								),
							},
							"virtual_host": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 200),
									validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-\/*:_+=.@-]*$`), ""),
								),
							},
						},
					},
				},
				"self_managed_kafka": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.sqs_queue_parameters",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "100"
								},
							},
							"consumer_group_id": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 200),
									validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9-\/*:_+=.@-]*$`), ""),
								),
							},
							"credentials": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"basic_auth": {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
										"client_certificate_tls_auth": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
										"sasl_scram_256_auth": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
										"sasl_scram_512_auth": {
											Type:         schema.TypeString,
											Optional:     true,
											ValidateFunc: verifySecretsManagerARN,
										},
									},
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
							"server_root_ca_certificate": {
								Type:         schema.TypeString,
								Optional:     true,
								ValidateFunc: verify.ValidARN,
							},
							"servers": {
								Type:     schema.TypeSet,
								Optional: true,
								ForceNew: true,
								MaxItems: 2,
								Elem: &schema.Schema{
									Type: schema.TypeString,
									ValidateFunc: validation.All(
										validation.StringLenBetween(1, 300),
										validation.StringMatch(regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9]):[0-9]{1,5}$`), ""),
									),
								},
							},
							"starting_position": {
								Type:             schema.TypeString,
								Optional:         true,
								ForceNew:         true,
								ValidateDiagFunc: enum.Validate[types.SelfManagedKafkaStartPosition](),
							},
							"topic": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
								ValidateFunc: validation.All(
									validation.StringLenBetween(1, 249),
									validation.StringMatch(regexp.MustCompile(`^[^.]([a-zA-Z0-9\-_.]+)$`), ""),
								),
							},
							"vpc": {
								Type:     schema.TypeList,
								Optional: true,
								MaxItems: 1,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"security_groups": {
											Type:     schema.TypeSet,
											Optional: true,
											MaxItems: 5,
											Elem: &schema.Schema{
												Type: schema.TypeString,
												ValidateFunc: validation.All(
													validation.StringLenBetween(1, 1024),
													validation.StringMatch(regexp.MustCompile(`^sg-[0-9a-zA-Z]*$`), ""),
												),
											},
										},
										"subnets": {
											Type:     schema.TypeSet,
											Optional: true,
											MaxItems: 16,
											Elem: &schema.Schema{
												Type: schema.TypeString,
												ValidateFunc: validation.All(
													validation.StringLenBetween(1, 1024),
													validation.StringMatch(regexp.MustCompile(`^subnet-[0-9a-z]*$`), ""),
												),
											},
										},
									},
								},
							},
						},
					},
				},
				"sqs_queue_parameters": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					ConflictsWith: []string{
						"source_parameters.0.activemq_broker_parameters",
						"source_parameters.0.dynamodb_stream_parameters",
						"source_parameters.0.kinesis_stream_parameters",
						"source_parameters.0.managed_streaming_kafka_parameters",
						"source_parameters.0.rabbit_mq_broker",
						"source_parameters.0.self_managed_kafka",
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"batch_size": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 10000),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "10"
								},
							},
							"maximum_batching_window_in_seconds": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(0, 300),
								DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
									if new != "" && new != "0" {
										return false
									}
									return old == "0"
								},
							},
						},
					},
				},
			},
		},
	}
}

func expandPipeSourceParameters(tfMap map[string]interface{}) *types.PipeSourceParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceParameters{}

	if v, ok := tfMap["activemq_broker_parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.ActiveMQBrokerParameters = expandPipeSourceActiveMQBrokerParameters(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["dynamodb_stream_parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.DynamoDBStreamParameters = expandPipeSourceDynamoDBStreamParameters(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["filter_criteria"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.FilterCriteria = expandFilterCriteria(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["kinesis_stream_parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.KinesisStreamParameters = expandPipeSourceKinesisStreamParameters(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["managed_streaming_kafka_parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.ManagedStreamingKafkaParameters = expandPipeSourceManagedStreamingKafkaParameters(v[0].(map[string]interface{}))
	}

	// TODO

	if v, ok := tfMap["sqs_queue_parameters"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.SqsQueueParameters = expandPipeSourceSqsQueueParameters(v[0].(map[string]interface{}))
	}

	return apiObject
}

func expandUpdatePipeSourceParameters(tfMap map[string]interface{}) *types.UpdatePipeSourceParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.UpdatePipeSourceParameters{}

	// TODO

	return apiObject
}

func expandFilterCriteria(tfMap map[string]interface{}) *types.FilterCriteria {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.FilterCriteria{}

	if v, ok := tfMap["filter"].([]interface{}); ok && len(v) > 0 {
		apiObject.Filters = expandFilters(v)
	}

	return apiObject
}

func expandFilter(tfMap map[string]interface{}) *types.Filter {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.Filter{}

	if v, ok := tfMap["pattern"].(string); ok && v != "" {
		apiObject.Pattern = aws.String(v)
	}

	return apiObject
}

func expandFilters(tfList []interface{}) []types.Filter {
	if len(tfList) == 0 {
		return nil
	}

	var apiObjects []types.Filter

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		apiObject := expandFilter(tfMap)

		if apiObject == nil {
			continue
		}

		apiObjects = append(apiObjects, *apiObject)
	}

	return apiObjects
}

func expandPipeSourceActiveMQBrokerParameters(tfMap map[string]interface{}) *types.PipeSourceActiveMQBrokerParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceActiveMQBrokerParameters{}

	if v, ok := tfMap["batch_size"].(int); ok && v != 0 {
		apiObject.BatchSize = aws.Int32(int32(v))
	}

	if v, ok := tfMap["credentials"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.Credentials = expandMQBrokerAccessCredentialsMemberBasicAuth(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["maximum_batching_window_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumBatchingWindowInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["queue_name"].(string); ok && v != "" {
		apiObject.QueueName = aws.String(v)
	}

	return apiObject
}

func expandMQBrokerAccessCredentialsMemberBasicAuth(tfMap map[string]interface{}) types.MQBrokerAccessCredentials {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.MQBrokerAccessCredentialsMemberBasicAuth{}

	if v, ok := tfMap["basic_auth"].(string); ok && v != "" {
		apiObject.Value = v
	}

	return apiObject
}

func expandPipeSourceDynamoDBStreamParameters(tfMap map[string]interface{}) *types.PipeSourceDynamoDBStreamParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceDynamoDBStreamParameters{}

	if v, ok := tfMap["batch_size"].(int); ok && v != 0 {
		apiObject.BatchSize = aws.Int32(int32(v))
	}

	if v, ok := tfMap["dead_letter_config"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.DeadLetterConfig = expandDeadLetterConfig(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["maximum_batching_window_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumBatchingWindowInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["maximum_record_age_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumRecordAgeInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["maximum_retry_attempts"].(int); ok && v != 0 {
		apiObject.MaximumRetryAttempts = aws.Int32(int32(v))
	}

	if v, ok := tfMap["on_partial_batch_item_failure"].(string); ok && v != "" {
		apiObject.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(v)
	}

	if v, ok := tfMap["parallelization_factor"].(int); ok && v != 0 {
		apiObject.ParallelizationFactor = aws.Int32(int32(v))
	}

	if v, ok := tfMap["starting_position"].(string); ok && v != "" {
		apiObject.StartingPosition = types.DynamoDBStreamStartPosition(v)
	}

	return apiObject
}

func expandPipeSourceKinesisStreamParameters(tfMap map[string]interface{}) *types.PipeSourceKinesisStreamParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceKinesisStreamParameters{}

	if v, ok := tfMap["batch_size"].(int); ok && v != 0 {
		apiObject.BatchSize = aws.Int32(int32(v))
	}

	if v, ok := tfMap["dead_letter_config"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.DeadLetterConfig = expandDeadLetterConfig(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["maximum_batching_window_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumBatchingWindowInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["maximum_record_age_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumRecordAgeInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["maximum_retry_attempts"].(int); ok && v != 0 {
		apiObject.MaximumRetryAttempts = aws.Int32(int32(v))
	}

	if v, ok := tfMap["on_partial_batch_item_failure"].(string); ok && v != "" {
		apiObject.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(v)
	}

	if v, ok := tfMap["parallelization_factor"].(int); ok && v != 0 {
		apiObject.ParallelizationFactor = aws.Int32(int32(v))
	}

	if v, ok := tfMap["starting_position"].(string); ok && v != "" {
		apiObject.StartingPosition = types.KinesisStreamStartPosition(v)
	}

	if v, ok := tfMap["starting_position_timestamp"].(string); ok && v != "" {
		v, _ := time.Parse(time.RFC3339, v)

		apiObject.StartingPositionTimestamp = aws.Time(v)
	}

	return apiObject
}

func expandPipeSourceManagedStreamingKafkaParameters(tfMap map[string]interface{}) *types.PipeSourceManagedStreamingKafkaParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceManagedStreamingKafkaParameters{}

	if v, ok := tfMap["batch_size"].(int); ok && v != 0 {
		apiObject.BatchSize = aws.Int32(int32(v))
	}

	if v, ok := tfMap["consumer_group_id"].(string); ok && v != "" {
		apiObject.ConsumerGroupID = aws.String(v)
	}

	if v, ok := tfMap["credentials"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		apiObject.Credentials = expandMSKAccessCredentials(v[0].(map[string]interface{}))
	}

	if v, ok := tfMap["maximum_batching_window_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumBatchingWindowInSeconds = aws.Int32(int32(v))
	}

	if v, ok := tfMap["starting_position"].(string); ok && v != "" {
		apiObject.StartingPosition = types.MSKStartPosition(v)
	}

	if v, ok := tfMap["topic_name"].(string); ok && v != "" {
		apiObject.TopicName = aws.String(v)
	}

	return apiObject
}

func expandMSKAccessCredentials(tfMap map[string]interface{}) types.MSKAccessCredentials {
	if tfMap == nil {
		return nil
	}

	if v, ok := tfMap["client_certificate_tls_auth"].(string); ok && v != "" {
		apiObject := &types.MSKAccessCredentialsMemberClientCertificateTlsAuth{
			Value: v,
		}

		return apiObject
	}

	if v, ok := tfMap["sasl_scram_512_auth"].(string); ok && v != "" {
		apiObject := &types.MSKAccessCredentialsMemberSaslScram512Auth{
			Value: v,
		}

		return apiObject
	}

	return nil
}

func expandPipeSourceSqsQueueParameters(tfMap map[string]interface{}) *types.PipeSourceSqsQueueParameters {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.PipeSourceSqsQueueParameters{}

	if v, ok := tfMap["batch_size"].(int); ok && v != 0 {
		apiObject.BatchSize = aws.Int32(int32(v))
	}

	if v, ok := tfMap["maximum_batching_window_in_seconds"].(int); ok && v != 0 {
		apiObject.MaximumBatchingWindowInSeconds = aws.Int32(int32(v))
	}

	return apiObject
}

func expandDeadLetterConfig(tfMap map[string]interface{}) *types.DeadLetterConfig {
	if tfMap == nil {
		return nil
	}

	apiObject := &types.DeadLetterConfig{}

	if v, ok := tfMap["arn"].(string); ok && v != "" {
		apiObject.Arn = aws.String(v)
	}

	return apiObject
}

func flattenPipeSourceParameters(apiObject *types.PipeSourceParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.ActiveMQBrokerParameters; v != nil {
		tfMap["activemq_broker_parameters"] = []interface{}{flattenPipeSourceActiveMQBrokerParameters(v)}
	}

	if v := apiObject.DynamoDBStreamParameters; v != nil {
		tfMap["dynamodb_stream_parameters"] = []interface{}{flattenPipeSourceDynamoDBStreamParameters(v)}
	}

	if v := apiObject.FilterCriteria; v != nil {
		tfMap["filter_criteria"] = []interface{}{flattenFilterCriteria(v)}
	}

	if v := apiObject.KinesisStreamParameters; v != nil {
		tfMap["kinesis_stream_parameters"] = []interface{}{flattenPipeSourceKinesisStreamParameters(v)}
	}

	if v := apiObject.ManagedStreamingKafkaParameters; v != nil {
		tfMap["managed_streaming_kafka_parameters"] = []interface{}{flattenPipeSourceManagedStreamingKafkaParameters(v)}
	}

	// TODO

	if v := apiObject.SqsQueueParameters; v != nil {
		tfMap["sqs_queue_parameters"] = []interface{}{flattenPipeSourceSqsQueueParameters(v)}
	}

	return tfMap
}

func flattenFilterCriteria(apiObject *types.FilterCriteria) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Filters; v != nil {
		tfMap["filter"] = flattenFilters(v)
	}

	return tfMap
}

func flattenFilter(apiObject types.Filter) map[string]interface{} {
	tfMap := map[string]interface{}{}

	if v := apiObject.Pattern; v != nil {
		tfMap["pattern"] = aws.ToString(v)
	}

	return tfMap
}

func flattenFilters(apiObjects []types.Filter) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfList = append(tfList, flattenFilter(apiObject))
	}

	return tfList
}

func flattenPipeSourceActiveMQBrokerParameters(apiObject *types.PipeSourceActiveMQBrokerParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BatchSize; v != nil {
		tfMap["batch_size"] = aws.ToInt32(v)
	}

	if v := apiObject.Credentials; v != nil {
		tfMap["credentials"] = []interface{}{flattenMQBrokerAccessCredentialsMemberBasicAuth(v.(*types.MQBrokerAccessCredentialsMemberBasicAuth))}
	}

	if v := apiObject.MaximumBatchingWindowInSeconds; v != nil {
		tfMap["maximum_batching_window_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.QueueName; v != nil {
		tfMap["queue_name"] = aws.ToString(v)
	}

	return tfMap
}

func flattenMQBrokerAccessCredentialsMemberBasicAuth(apiObject *types.MQBrokerAccessCredentialsMemberBasicAuth) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Value; v != "" {
		tfMap["basic_auth"] = v
	}

	return tfMap
}

func flattenPipeSourceDynamoDBStreamParameters(apiObject *types.PipeSourceDynamoDBStreamParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BatchSize; v != nil {
		tfMap["batch_size"] = aws.ToInt32(v)
	}

	if v := apiObject.DeadLetterConfig; v != nil {
		tfMap["dead_letter_config"] = []interface{}{flattenDeadLetterConfig(v)}
	}

	if v := apiObject.MaximumBatchingWindowInSeconds; v != nil {
		tfMap["maximum_batching_window_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.MaximumRecordAgeInSeconds; v != nil {
		tfMap["maximum_record_age_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.MaximumRetryAttempts; v != nil {
		tfMap["maximum_retry_attempts"] = aws.ToInt32(v)
	}

	if v := apiObject.OnPartialBatchItemFailure; v != "" {
		tfMap["on_partial_batch_item_failure"] = v
	}

	if v := apiObject.ParallelizationFactor; v != nil {
		tfMap["parallelization_factor"] = aws.ToInt32(v)
	}

	if v := apiObject.StartingPosition; v != "" {
		tfMap["starting_position"] = v
	}

	return tfMap
}

func flattenPipeSourceKinesisStreamParameters(apiObject *types.PipeSourceKinesisStreamParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BatchSize; v != nil {
		tfMap["batch_size"] = aws.ToInt32(v)
	}

	if v := apiObject.DeadLetterConfig; v != nil {
		tfMap["dead_letter_config"] = []interface{}{flattenDeadLetterConfig(v)}
	}

	if v := apiObject.MaximumBatchingWindowInSeconds; v != nil {
		tfMap["maximum_batching_window_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.MaximumRecordAgeInSeconds; v != nil {
		tfMap["maximum_record_age_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.MaximumRetryAttempts; v != nil {
		tfMap["maximum_retry_attempts"] = aws.ToInt32(v)
	}

	if v := apiObject.OnPartialBatchItemFailure; v != "" {
		tfMap["on_partial_batch_item_failure"] = v
	}

	if v := apiObject.ParallelizationFactor; v != nil {
		tfMap["parallelization_factor"] = aws.ToInt32(v)
	}

	if v := apiObject.StartingPosition; v != "" {
		tfMap["starting_position"] = v
	}

	if v := apiObject.StartingPositionTimestamp; v != nil {
		tfMap["starting_position_timestamp"] = aws.ToTime(v).Format(time.RFC3339)
	}

	return tfMap
}

func flattenPipeSourceManagedStreamingKafkaParameters(apiObject *types.PipeSourceManagedStreamingKafkaParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BatchSize; v != nil {
		tfMap["batch_size"] = aws.ToInt32(v)
	}

	if v := apiObject.ConsumerGroupID; v != nil {
		tfMap["consumer_group_id"] = aws.ToString(v)
	}

	if v := apiObject.Credentials; v != nil {
		tfMap["credentials"] = []interface{}{flattenMSKAccessCredentials(v)}
	}

	if v := apiObject.MaximumBatchingWindowInSeconds; v != nil {
		tfMap["maximum_batching_window_in_seconds"] = aws.ToInt32(v)
	}

	if v := apiObject.StartingPosition; v != "" {
		tfMap["starting_position"] = v
	}

	if v := apiObject.TopicName; v != nil {
		tfMap["topic_name"] = aws.ToString(v)
	}

	return tfMap
}

func flattenMSKAccessCredentials(apiObject types.MSKAccessCredentials) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if apiObject, ok := apiObject.(*types.MSKAccessCredentialsMemberClientCertificateTlsAuth); ok {
		if v := apiObject.Value; v != "" {
			tfMap["client_certificate_tls_auth"] = v
		}
	}

	if apiObject, ok := apiObject.(*types.MSKAccessCredentialsMemberSaslScram512Auth); ok {
		if v := apiObject.Value; v != "" {
			tfMap["sasl_scram_512_auth"] = v
		}
	}

	return tfMap
}

func flattenPipeSourceSqsQueueParameters(apiObject *types.PipeSourceSqsQueueParameters) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.BatchSize; v != nil {
		tfMap["batch_size"] = aws.ToInt32(v)
	}

	if v := apiObject.MaximumBatchingWindowInSeconds; v != nil {
		tfMap["maximum_batching_window_in_seconds"] = aws.ToInt32(v)
	}

	return tfMap
}

func flattenDeadLetterConfig(apiObject *types.DeadLetterConfig) map[string]interface{} {
	if apiObject == nil {
		return nil
	}

	tfMap := map[string]interface{}{}

	if v := apiObject.Arn; v != nil {
		tfMap["arn"] = aws.ToString(v)
	}

	return tfMap
}

/*
func expandSourceParameters(config []interface{}) *types.PipeSourceParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceParameters
	for _, c := range config {
		param, ok := c.(map[string]interface{})
		if !ok {
			return nil
		}

		if val, ok := param["active_mq_broker"]; ok {
			parameters.ActiveMQBrokerParameters = expandSourceActiveMQBrokerParameters(val.([]interface{}))
		}

		if val, ok := param["dynamo_db_stream"]; ok {
			parameters.DynamoDBStreamParameters = expandSourceDynamoDBStreamParameters(val.([]interface{}))
		}

		if val, ok := param["kinesis_stream"]; ok {
			parameters.KinesisStreamParameters = expandSourceKinesisStreamParameters(val.([]interface{}))
		}

		if val, ok := param["managed_streaming_kafka"]; ok {
			parameters.ManagedStreamingKafkaParameters = expandSourceManagedStreamingKafkaParameters(val.([]interface{}))
		}

		if val, ok := param["rabbit_mq_broker"]; ok {
			parameters.RabbitMQBrokerParameters = expandSourceRabbitMQBrokerParameters(val.([]interface{}))
		}

		if val, ok := param["self_managed_kafka"]; ok {
			parameters.SelfManagedKafkaParameters = expandSourceSelfManagedKafkaParameters(val.([]interface{}))
		}

		if val, ok := param["sqs_queue"]; ok {
			parameters.SqsQueueParameters = expandSourceSqsQueueParameters(val.([]interface{}))
		}

		if val, ok := param["filter_criteria"]; ok {
			parameters.FilterCriteria = expandSourceFilterCriteria(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceActiveMQBrokerParameters(config []interface{}) *types.PipeSourceActiveMQBrokerParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceActiveMQBrokerParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.QueueName = expandString("queue", param)
		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				var credentialsParameters types.MQBrokerAccessCredentialsMemberBasicAuth
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
				}
				parameters.Credentials = &credentialsParameters
			}
		}
	}
	return &parameters
}

func expandSourceDynamoDBStreamParameters(config []interface{}) *types.PipeSourceDynamoDBStreamParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceDynamoDBStreamParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.MaximumRecordAgeInSeconds = expandInt32("maximum_record_age_in_seconds", param)
		parameters.ParallelizationFactor = expandInt32("parallelization_factor", param)
		parameters.MaximumRetryAttempts = expandInt32("maximum_retry_attempts", param)
		startingPosition := expandStringValue("starting_position", param)
		if startingPosition != "" {
			parameters.StartingPosition = types.DynamoDBStreamStartPosition(startingPosition)
		}
		onPartialBatchItemFailure := expandStringValue("on_partial_batch_item_failure", param)
		if onPartialBatchItemFailure != "" {
			parameters.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(onPartialBatchItemFailure)
		}
		if val, ok := param["dead_letter_config"]; ok {
			parameters.DeadLetterConfig = expandSourceDeadLetterConfig(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceKinesisStreamParameters(config []interface{}) *types.PipeSourceKinesisStreamParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceKinesisStreamParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.MaximumRecordAgeInSeconds = expandInt32("maximum_record_age_in_seconds", param)
		parameters.ParallelizationFactor = expandInt32("parallelization_factor", param)
		parameters.MaximumRetryAttempts = expandInt32("maximum_retry_attempts", param)

		startingPosition := expandStringValue("starting_position", param)
		if startingPosition != "" {
			parameters.StartingPosition = types.KinesisStreamStartPosition(startingPosition)
		}
		onPartialBatchItemFailure := expandStringValue("on_partial_batch_item_failure", param)
		if onPartialBatchItemFailure != "" {
			parameters.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(onPartialBatchItemFailure)
		}
		if val, ok := param["starting_position_timestamp"]; ok {
			t, _ := time.Parse(time.RFC3339, val.(string))

			parameters.StartingPositionTimestamp = aws.Time(t)
		}
		if val, ok := param["dead_letter_config"]; ok {
			parameters.DeadLetterConfig = expandSourceDeadLetterConfig(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceManagedStreamingKafkaParameters(config []interface{}) *types.PipeSourceManagedStreamingKafkaParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceManagedStreamingKafkaParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.TopicName = expandString("topic", param)
		parameters.ConsumerGroupID = expandString("consumer_group_id", param)

		startingPosition := expandStringValue("starting_position", param)
		if startingPosition != "" {
			parameters.StartingPosition = types.MSKStartPosition(startingPosition)
		}

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					if _, ok := credentialsParam["client_certificate_tls_auth"]; ok {
						var credentialsParameters types.MSKAccessCredentialsMemberClientCertificateTlsAuth
						credentialsParameters.Value = expandStringValue("client_certificate_tls_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_512_auth"]; ok {
						var credentialsParameters types.MSKAccessCredentialsMemberSaslScram512Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_512_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
				}
			}
		}
	}
	return &parameters
}

func expandSourceRabbitMQBrokerParameters(config []interface{}) *types.PipeSourceRabbitMQBrokerParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceRabbitMQBrokerParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.QueueName = expandString("queue", param)
		parameters.VirtualHost = expandString("virtual_host", param)

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				var credentialsParameters types.MQBrokerAccessCredentialsMemberBasicAuth
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
				}
				parameters.Credentials = &credentialsParameters
			}
		}
	}
	return &parameters
}

func expandSourceSelfManagedKafkaParameters(config []interface{}) *types.PipeSourceSelfManagedKafkaParameters {
	if len(config) == 0 {
		return nil
	}
	var parameters types.PipeSourceSelfManagedKafkaParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.TopicName = expandString("topic", param)
		parameters.ConsumerGroupID = expandString("consumer_group_id", param)
		parameters.ServerRootCaCertificate = expandString("server_root_ca_certificate", param)
		startingPosition := expandStringValue("starting_position", param)
		if startingPosition != "" {
			parameters.StartingPosition = types.SelfManagedKafkaStartPosition(startingPosition)
		}
		if value, ok := param["servers"]; ok && value.(*schema.Set).Len() > 0 {
			parameters.AdditionalBootstrapServers = flex.ExpandStringValueSet(value.(*schema.Set))
		}

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					if _, ok := credentialsParam["basic_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberBasicAuth
						credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["client_certificate_tls_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberClientCertificateTlsAuth
						credentialsParameters.Value = expandStringValue("client_certificate_tls_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_512_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram512Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_512_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_256_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram256Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_256_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
				}
			}
		}

		if val, ok := param["vpc"]; ok {
			vpcConfig := val.([]interface{})
			if len(vpcConfig) != 0 {
				var vpcParameters types.SelfManagedKafkaAccessConfigurationVpc
				for _, vc := range vpcConfig {
					vpcParam := vc.(map[string]interface{})
					if value, ok := vpcParam["security_groups"]; ok && value.(*schema.Set).Len() > 0 {
						vpcParameters.SecurityGroup = flex.ExpandStringValueSet(value.(*schema.Set))
					}
					if value, ok := vpcParam["subnets"]; ok && value.(*schema.Set).Len() > 0 {
						vpcParameters.Subnets = flex.ExpandStringValueSet(value.(*schema.Set))
					}
				}
				parameters.Vpc = &vpcParameters
			}
		}
	}

	return &parameters
}

func expandSourceSqsQueueParameters(config []interface{}) *types.PipeSourceSqsQueueParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.PipeSourceSqsQueueParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
	}

	return &parameters
}

func expandSourceDeadLetterConfig(config []interface{}) *types.DeadLetterConfig {
	if len(config) == 0 {
		return nil
	}

	var parameters types.DeadLetterConfig
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.Arn = expandString("arn", param)
	}

	return &parameters
}

func expandSourceFilterCriteria(config []interface{}) *types.FilterCriteria {
	if len(config) == 0 {
		return nil
	}

	var parameters types.FilterCriteria
	for _, c := range config {
		param := c.(map[string]interface{})
		if val, ok := param["filter"]; ok {
			filtersConfig := val.([]interface{})
			var filters []types.Filter
			for _, f := range filtersConfig {
				filterParam := f.(map[string]interface{})
				pattern := expandString("pattern", filterParam)
				if pattern != nil {
					filters = append(filters, types.Filter{
						Pattern: pattern,
					})
				}
			}
			if len(filters) > 0 {
				parameters.Filters = filters
			}
		}
	}

	return &parameters
}

func flattenSourceParameters(sourceParameters *types.PipeSourceParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if sourceParameters.ActiveMQBrokerParameters != nil {
		config["active_mq_broker"] = flattenSourceActiveMQBrokerParameters(sourceParameters.ActiveMQBrokerParameters)
	}

	if sourceParameters.DynamoDBStreamParameters != nil {
		config["dynamo_db_stream"] = flattenSourceDynamoDBStreamParameters(sourceParameters.DynamoDBStreamParameters)
	}

	if sourceParameters.KinesisStreamParameters != nil {
		config["kinesis_stream"] = flattenSourceKinesisStreamParameters(sourceParameters.KinesisStreamParameters)
	}

	if sourceParameters.ManagedStreamingKafkaParameters != nil {
		config["managed_streaming_kafka"] = flattenSourceManagedStreamingKafkaParameters(sourceParameters.ManagedStreamingKafkaParameters)
	}

	if sourceParameters.RabbitMQBrokerParameters != nil {
		config["rabbit_mq_broker"] = flattenSourceRabbitMQBrokerParameters(sourceParameters.RabbitMQBrokerParameters)
	}

	if sourceParameters.SelfManagedKafkaParameters != nil {
		config["self_managed_kafka"] = flattenSourceSelfManagedKafkaParameters(sourceParameters.SelfManagedKafkaParameters)
	}

	if sourceParameters.SqsQueueParameters != nil {
		config["sqs_queue"] = flattenSourceSqsQueueParameters(sourceParameters.SqsQueueParameters)
	}

	if sourceParameters.FilterCriteria != nil {
		criteria := flattenSourceFilterCriteria(sourceParameters.FilterCriteria)
		if len(criteria) > 0 {
			config["filter_criteria"] = criteria
		}
	}

	if len(config) == 0 {
		return nil
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceActiveMQBrokerParameters(parameters *types.PipeSourceActiveMQBrokerParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.QueueName != nil {
		config["queue"] = aws.ToString(parameters.QueueName)
	}
	if parameters.Credentials != nil {
		credentialsConfig := make(map[string]interface{})
		switch v := parameters.Credentials.(type) {
		case *types.MQBrokerAccessCredentialsMemberBasicAuth:
			credentialsConfig["basic_auth"] = v.Value
		}
		config["credentials"] = []map[string]interface{}{credentialsConfig}
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceDynamoDBStreamParameters(parameters *types.PipeSourceDynamoDBStreamParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.MaximumRecordAgeInSeconds != nil {
		config["maximum_record_age_in_seconds"] = aws.ToInt32(parameters.MaximumRecordAgeInSeconds)
	}
	if parameters.ParallelizationFactor != nil {
		config["parallelization_factor"] = aws.ToInt32(parameters.ParallelizationFactor)
	}
	if parameters.MaximumRetryAttempts != nil {
		config["maximum_retry_attempts"] = aws.ToInt32(parameters.MaximumRetryAttempts)
	}
	if parameters.StartingPosition != "" {
		config["starting_position"] = parameters.StartingPosition
	}
	if parameters.OnPartialBatchItemFailure != "" {
		config["on_partial_batch_item_failure"] = parameters.OnPartialBatchItemFailure
	}
	if parameters.DeadLetterConfig != nil {
		config["dead_letter_config"] = flattenSourceDeadLetterConfig(parameters.DeadLetterConfig)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceKinesisStreamParameters(parameters *types.PipeSourceKinesisStreamParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.MaximumRecordAgeInSeconds != nil {
		config["maximum_record_age_in_seconds"] = aws.ToInt32(parameters.MaximumRecordAgeInSeconds)
	}
	if parameters.ParallelizationFactor != nil {
		config["parallelization_factor"] = aws.ToInt32(parameters.ParallelizationFactor)
	}
	if parameters.MaximumRetryAttempts != nil {
		config["maximum_retry_attempts"] = aws.ToInt32(parameters.MaximumRetryAttempts)
	}
	if parameters.StartingPosition != "" {
		config["starting_position"] = parameters.StartingPosition
	}
	if parameters.OnPartialBatchItemFailure != "" {
		config["on_partial_batch_item_failure"] = parameters.OnPartialBatchItemFailure
	}
	if parameters.StartingPositionTimestamp != nil {
		config["starting_position_timestamp"] = aws.ToTime(parameters.StartingPositionTimestamp).Format(time.RFC3339)
	}
	if parameters.DeadLetterConfig != nil {
		config["dead_letter_config"] = flattenSourceDeadLetterConfig(parameters.DeadLetterConfig)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceManagedStreamingKafkaParameters(parameters *types.PipeSourceManagedStreamingKafkaParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.ConsumerGroupID != nil {
		config["consumer_group_id"] = aws.ToString(parameters.ConsumerGroupID)
	}
	if parameters.StartingPosition != "" {
		config["starting_position"] = parameters.StartingPosition
	}
	if parameters.TopicName != nil {
		config["topic"] = aws.ToString(parameters.TopicName)
	}
	if parameters.Credentials != nil {
		credentialsConfig := make(map[string]interface{})
		switch v := parameters.Credentials.(type) {
		case *types.MSKAccessCredentialsMemberClientCertificateTlsAuth:
			credentialsConfig["client_certificate_tls_auth"] = v.Value
		case *types.MSKAccessCredentialsMemberSaslScram512Auth:
			credentialsConfig["sasl_scram_512_auth"] = v.Value
		}
		config["credentials"] = []map[string]interface{}{credentialsConfig}
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceRabbitMQBrokerParameters(parameters *types.PipeSourceRabbitMQBrokerParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.QueueName != nil {
		config["queue"] = aws.ToString(parameters.QueueName)
	}
	if parameters.VirtualHost != nil {
		config["virtual_host"] = aws.ToString(parameters.VirtualHost)
	}
	if parameters.Credentials != nil {
		credentialsConfig := make(map[string]interface{})
		switch v := parameters.Credentials.(type) {
		case *types.MQBrokerAccessCredentialsMemberBasicAuth:
			credentialsConfig["basic_auth"] = v.Value
		}
		config["credentials"] = []map[string]interface{}{credentialsConfig}
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceSelfManagedKafkaParameters(parameters *types.PipeSourceSelfManagedKafkaParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}
	if parameters.ConsumerGroupID != nil {
		config["consumer_group_id"] = aws.ToString(parameters.ConsumerGroupID)
	}
	if parameters.StartingPosition != "" {
		config["starting_position"] = parameters.StartingPosition
	}
	if parameters.TopicName != nil {
		config["topic"] = aws.ToString(parameters.TopicName)
	}
	if parameters.AdditionalBootstrapServers != nil {
		config["servers"] = flex.FlattenStringValueSet(parameters.AdditionalBootstrapServers)
	}
	if parameters.ServerRootCaCertificate != nil {
		config["server_root_ca_certificate"] = aws.ToString(parameters.ServerRootCaCertificate)
	}

	if parameters.Credentials != nil {
		credentialsConfig := make(map[string]interface{})
		switch v := parameters.Credentials.(type) {
		case *types.SelfManagedKafkaAccessConfigurationCredentialsMemberBasicAuth:
			credentialsConfig["basic_auth"] = v.Value
		case *types.SelfManagedKafkaAccessConfigurationCredentialsMemberClientCertificateTlsAuth:
			credentialsConfig["client_certificate_tls_auth"] = v.Value
		case *types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram256Auth:
			credentialsConfig["sasl_scram_256_auth"] = v.Value
		case *types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram512Auth:
			credentialsConfig["sasl_scram_512_auth"] = v.Value
		}
		config["credentials"] = []map[string]interface{}{credentialsConfig}
	}
	if parameters.Vpc != nil {
		vpcConfig := make(map[string]interface{})
		vpcConfig["security_groups"] = flex.FlattenStringValueSet(parameters.Vpc.SecurityGroup)
		vpcConfig["subnets"] = flex.FlattenStringValueSet(parameters.Vpc.Subnets)
		config["vpc"] = []map[string]interface{}{vpcConfig}
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceSqsQueueParameters(parameters *types.PipeSourceSqsQueueParameters) []map[string]interface{} {
	config := make(map[string]interface{})

	if parameters.BatchSize != nil {
		config["batch_size"] = aws.ToInt32(parameters.BatchSize)
	}
	if parameters.MaximumBatchingWindowInSeconds != nil && aws.ToInt32(parameters.MaximumBatchingWindowInSeconds) != 0 {
		config["maximum_batching_window_in_seconds"] = aws.ToInt32(parameters.MaximumBatchingWindowInSeconds)
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceFilterCriteria(parameters *types.FilterCriteria) []map[string]interface{} {
	config := make(map[string]interface{})

	if len(parameters.Filters) != 0 {
		var filters []map[string]interface{}
		for _, filter := range parameters.Filters {
			pattern := make(map[string]interface{})
			pattern["pattern"] = aws.ToString(filter.Pattern)
			filters = append(filters, pattern)
		}
		if len(filters) != 0 {
			config["filter"] = filters
		}
	}

	result := []map[string]interface{}{config}
	return result
}

func flattenSourceDeadLetterConfig(parameters *types.DeadLetterConfig) []map[string]interface{} {
	if parameters == nil {
		return nil
	}

	config := make(map[string]interface{})
	if parameters.Arn != nil {
		config["arn"] = aws.ToString(parameters.Arn)
	}

	result := []map[string]interface{}{config}
	return result
}

func expandSourceUpdateParameters(config []interface{}) *types.UpdatePipeSourceParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceParameters
	for _, c := range config {
		param, ok := c.(map[string]interface{})
		if !ok {
			return nil
		}

		if val, ok := param["active_mq_broker"]; ok {
			parameters.ActiveMQBrokerParameters = expandSourceUpdateActiveMQBrokerParameters(val.([]interface{}))
		}

		if val, ok := param["dynamo_db_stream"]; ok {
			parameters.DynamoDBStreamParameters = expandSourceUpdateDynamoDBStreamParameters(val.([]interface{}))
		}

		if val, ok := param["kinesis_stream"]; ok {
			parameters.KinesisStreamParameters = expandSourceUpdateKinesisStreamParameters(val.([]interface{}))
		}

		if val, ok := param["managed_streaming_kafka"]; ok {
			parameters.ManagedStreamingKafkaParameters = expandSourceUpdateManagedStreamingKafkaParameters(val.([]interface{}))
		}

		if val, ok := param["rabbit_mq_broker"]; ok {
			parameters.RabbitMQBrokerParameters = expandSourceUpdateRabbitMQBrokerParameters(val.([]interface{}))
		}

		if val, ok := param["self_managed_kafka"]; ok {
			parameters.SelfManagedKafkaParameters = expandSourceUpdateSelfManagedKafkaParameters(val.([]interface{}))
		}

		if val, ok := param["sqs_queue"]; ok {
			parameters.SqsQueueParameters = expandSourceUpdateSqsQueueParameters(val.([]interface{}))
		}

		if val, ok := param["filter_criteria"]; ok {
			parameters.FilterCriteria = expandSourceFilterCriteria(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceUpdateActiveMQBrokerParameters(config []interface{}) *types.UpdatePipeSourceActiveMQBrokerParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceActiveMQBrokerParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				var credentialsParameters types.MQBrokerAccessCredentialsMemberBasicAuth
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
				}
				parameters.Credentials = &credentialsParameters
			}
		}
	}
	return &parameters
}

func expandSourceUpdateDynamoDBStreamParameters(config []interface{}) *types.UpdatePipeSourceDynamoDBStreamParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceDynamoDBStreamParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.MaximumRecordAgeInSeconds = expandInt32("maximum_record_age_in_seconds", param)
		parameters.ParallelizationFactor = expandInt32("parallelization_factor", param)
		parameters.MaximumRetryAttempts = expandInt32("maximum_retry_attempts", param)
		onPartialBatchItemFailure := expandStringValue("on_partial_batch_item_failure", param)
		if onPartialBatchItemFailure != "" {
			parameters.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(onPartialBatchItemFailure)
		}
		if val, ok := param["dead_letter_config"]; ok {
			parameters.DeadLetterConfig = expandSourceDeadLetterConfig(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceUpdateKinesisStreamParameters(config []interface{}) *types.UpdatePipeSourceKinesisStreamParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceKinesisStreamParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.MaximumRecordAgeInSeconds = expandInt32("maximum_record_age_in_seconds", param)
		parameters.ParallelizationFactor = expandInt32("parallelization_factor", param)
		parameters.MaximumRetryAttempts = expandInt32("maximum_retry_attempts", param)

		onPartialBatchItemFailure := expandStringValue("on_partial_batch_item_failure", param)
		if onPartialBatchItemFailure != "" {
			parameters.OnPartialBatchItemFailure = types.OnPartialBatchItemFailureStreams(onPartialBatchItemFailure)
		}
		if val, ok := param["dead_letter_config"]; ok {
			parameters.DeadLetterConfig = expandSourceDeadLetterConfig(val.([]interface{}))
		}
	}
	return &parameters
}

func expandSourceUpdateManagedStreamingKafkaParameters(config []interface{}) *types.UpdatePipeSourceManagedStreamingKafkaParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceManagedStreamingKafkaParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					if _, ok := credentialsParam["client_certificate_tls_auth"]; ok {
						var credentialsParameters types.MSKAccessCredentialsMemberClientCertificateTlsAuth
						credentialsParameters.Value = expandStringValue("client_certificate_tls_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_512_auth"]; ok {
						var credentialsParameters types.MSKAccessCredentialsMemberSaslScram512Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_512_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
				}
			}
		}
	}
	return &parameters
}

func expandSourceUpdateRabbitMQBrokerParameters(config []interface{}) *types.UpdatePipeSourceRabbitMQBrokerParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceRabbitMQBrokerParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				var credentialsParameters types.MQBrokerAccessCredentialsMemberBasicAuth
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
				}
				parameters.Credentials = &credentialsParameters
			}
		}
	}
	return &parameters
}

func expandSourceUpdateSelfManagedKafkaParameters(config []interface{}) *types.UpdatePipeSourceSelfManagedKafkaParameters {
	if len(config) == 0 {
		return nil
	}
	var parameters types.UpdatePipeSourceSelfManagedKafkaParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
		parameters.ServerRootCaCertificate = expandString("server_root_ca_certificate", param)

		if val, ok := param["credentials"]; ok {
			credentialsConfig := val.([]interface{})
			if len(credentialsConfig) != 0 {
				for _, cc := range credentialsConfig {
					credentialsParam := cc.(map[string]interface{})
					if _, ok := credentialsParam["basic_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberBasicAuth
						credentialsParameters.Value = expandStringValue("basic_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["client_certificate_tls_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberClientCertificateTlsAuth
						credentialsParameters.Value = expandStringValue("client_certificate_tls_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_512_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram512Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_512_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
					if _, ok := credentialsParam["sasl_scram_256_auth"]; ok {
						var credentialsParameters types.SelfManagedKafkaAccessConfigurationCredentialsMemberSaslScram256Auth
						credentialsParameters.Value = expandStringValue("sasl_scram_256_auth", credentialsParam)
						parameters.Credentials = &credentialsParameters
					}
				}
			}
		}

		if val, ok := param["vpc"]; ok {
			vpcConfig := val.([]interface{})
			if len(vpcConfig) != 0 {
				var vpcParameters types.SelfManagedKafkaAccessConfigurationVpc
				for _, vc := range vpcConfig {
					vpcParam := vc.(map[string]interface{})
					if value, ok := vpcParam["security_groups"]; ok && value.(*schema.Set).Len() > 0 {
						vpcParameters.SecurityGroup = flex.ExpandStringValueSet(value.(*schema.Set))
					}
					if value, ok := vpcParam["subnets"]; ok && value.(*schema.Set).Len() > 0 {
						vpcParameters.Subnets = flex.ExpandStringValueSet(value.(*schema.Set))
					}
				}
				parameters.Vpc = &vpcParameters
			}
		}
	}

	return &parameters
}

func expandSourceUpdateSqsQueueParameters(config []interface{}) *types.UpdatePipeSourceSqsQueueParameters {
	if len(config) == 0 {
		return nil
	}

	var parameters types.UpdatePipeSourceSqsQueueParameters
	for _, c := range config {
		param := c.(map[string]interface{})
		parameters.BatchSize = expandInt32("batch_size", param)
		parameters.MaximumBatchingWindowInSeconds = expandInt32("maximum_batching_window_in_seconds", param)
	}

	return &parameters
}
*/
