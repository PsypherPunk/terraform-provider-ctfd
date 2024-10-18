package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	teams, err := client.GetTeams()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("teams", teams); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceTeams() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of the registered teams.",
		ReadContext: dataSourceTeamsRead,
		Schema: map[string]*schema.Schema{
			"teams": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email": {
							Type:     schema.TypeString,
							Required: true,
						},
						"website": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"affiliation": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"country": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"hidden": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"banned": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"captain_id": {
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"bracket": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"secret": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"oauth_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"members": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							MinItems: 0,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 0,
						},
					},
				},
			},
		},
	}
}
