package provider

import (
	"context"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChallengesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	challenges, err := client.GetChallenges()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("challenges", challenges); err != nil {
		return diag.FromErr(err)
	}

	// TODO: determine based on contents of `challenges`?
	d.SetId("challenges")

	return diags
}

func dataSourceChallenges() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of the current challenges.",
		ReadContext: dataSourceChallengesRead,
		Schema: map[string]*schema.Schema{
			"challenges": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solves": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solved_by_me": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"category": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"template": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"script": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
