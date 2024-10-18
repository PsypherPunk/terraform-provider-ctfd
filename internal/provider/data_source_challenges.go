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
			"challenges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solves": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solved_by_me": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"script": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
