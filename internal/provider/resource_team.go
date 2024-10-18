package provider

import (
	"context"
	"strconv"

	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	team := api.NewTeam{
		Name:        d.Get("name").(string),
		Email:       d.Get("email").(string),
		Password:    d.Get("password").(string),
		Website:     d.Get("website").(string),
		Affiliation: d.Get("affiliation").(string),
		Country:     d.Get("country").(string),
		Hidden:      d.Get("hidden").(bool),
		Banned:      d.Get("banned").(bool),
	}

	newTeam, err := client.CreateTeam(team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(newTeam.Id)))
	err = d.Set("members", newTeam.Members)
	if err != nil {
		return nil
	}
	err = d.Set("fields", newTeam.Fields)
	if err != nil {
		return nil
	}

	return diags
}

func resourceTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(string)

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	team, err := client.GetTeam(uint(intId))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(team.Id)))

	return diags
}

func resourceTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Id()

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	team := new(api.NewTeam)
	team.Password = d.Get("password").(string)
	team.Name = d.Get("name").(string)
	team.Email = d.Get("email").(string)
	team.Website = d.Get("website").(string)
	team.Affiliation = d.Get("affiliation").(string)
	team.Country = d.Get("country").(string)
	team.Hidden = d.Get("hidden").(bool)
	team.Banned = d.Get("banned").(bool)

	updatedTeam, err := client.UpdateTeam(uint(intId), *team)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(updatedTeam.Id)))

	return diags
}

func resourceTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	var diags diag.Diagnostics

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	err = client.DeleteTeam(uint(intId))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceTeam() *schema.Resource {
	return &schema.Resource{
		Description:   "Get details of a teams.",
		CreateContext: resourceTeamCreate,
		ReadContext:   resourceTeamRead,
		UpdateContext: resourceTeamUpdate,
		DeleteContext: resourceTeamDelete,
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
			"password": {
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
				Optional: true,
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
			},
		},
	}
}
