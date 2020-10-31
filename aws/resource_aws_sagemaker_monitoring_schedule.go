package aws

import (
	"errors"
	//"fmt"
	//"log"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sagemaker"
	//"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	//"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsSagemakerMonitoringSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSagemakerMonitoringScheduleCreate,
		Read:   resourceAwsSagemakerMonitoringScheduleRead,
		Update: resourceAwsSagemakerMonitoringScheduleUpdate,
		Delete: resourceAwsSagemakerMonitoringScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validateSagemakerName,
			},

			"monitoring_schedule_config": {
				Type:     schema.TypeList,
				Optional: false,
				MaxItems: 1,
				ForceNew: true, // todo: perhaps this can be an inplace update?
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule_expression": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},

						"monitoring_job_definition": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"baseline_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"constraints_resource": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_uri": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},

												"statistics_resource": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_uri": {
																Type:     schema.TypeString,
																Optional: true,
															},
														},
													},
												},
											},
										},
									},

									"monitoring_inputs": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint_input": {
													Type:     schema.TypeSet,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"endpoint_name": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validateSagemakerName,
															},

															"local_path": {
																Type:     schema.TypeString,
																Required: true,
															},

															"s3_data_distribution_type": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice(sagemaker.ProcessingS3DataDistributionType_Values(), false),
															},

															"s3_input_mode": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validation.StringInSlice(sagemaker.ProcessingS3InputMode_Values(), false),
															},
														},
													},
												},
											},
										},
									},

									"monitoring_output_config": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kms_key_id": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validateKmsKey,
												},

												"monitoring_outputs": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"s3_output": {
																Type:     schema.TypeSet,
																Required: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"s3_uri": {
																			Type:     schema.TypeString,
																			Required: true,
																		},

																		"s3_upload_mode": {
																			Type:         schema.TypeString,
																			Optional:     true,
																			ValidateFunc: validation.StringInSlice(sagemaker.ProcessingS3UploadMode_Values(), false),
																		},

																		"local_path": {
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

									"monitoring_resources": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_config": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"instance_count": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 100),
															},

															"instance_type": {
																Type:     schema.TypeString,
																Required: true,
															},

															"volume_size_in_gb": {
																Type:         schema.TypeInt,
																Required:     true,
																ValidateFunc: validation.IntBetween(1, 16384),
															},

															"volume_kms_key_id": {
																Type:         schema.TypeString,
																Optional:     true,
																ValidateFunc: validateKmsKey,
															},
														},
													},
												},
											},
										},
									},

									"monitoring_app_specification": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_uri": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validateSagemakerImage,
												},

												"container_entry_point": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 256,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},

												"container_arguments": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 256,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},

												"record_preprocessor_source_uri": {
													Type:     schema.TypeString,
													Optional: true,
												},

												"post_analytics_processor_source_uri": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},

									"stopping_condition": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"max_runtime_in_seconds": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntBetween(1, 86400),
												},
											},
										},
									},

									"environment": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"network_config": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"intercontainer_traffic_encryption": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"network_isolation": {
													Type:     schema.TypeBool,
													Optional: true,
												},

												"vpc_config": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													ForceNew: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"subnets": {
																Type:     schema.TypeSet,
																Required: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
																Set:      schema.HashString,
															},
															"security_group_ids": {
																Type:     schema.TypeSet,
																Required: true,
																Elem:     &schema.Schema{Type: schema.TypeString},
																Set:      schema.HashString,
															},
														},
													},
												},
											},
										},
									},

									"role_arn": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateArn,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceAwsSagemakerMonitoringScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("not implemented")
}

func resourceAwsSagemakerMonitoringScheduleRead(d *schema.ResourceData, meta interface{}) error {
	return errors.New("not implemented")
}

func resourceAwsSagemakerMonitoringScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	return errors.New("not implemented")
}

func resourceAwsSagemakerMonitoringScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	return errors.New("not implemented")
}
