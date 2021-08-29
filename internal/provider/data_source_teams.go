package provider

import (
	"context"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"strconv"
	"time"

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
			"teams": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"email": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"website": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"affiliation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"country": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"hidden": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"banned": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"captain_id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
							Optional: true,
						},
						"bracket": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"secret": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"oauth_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"members": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							MinItems: 0,
						},
						"created": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"fields": &schema.Schema{
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
