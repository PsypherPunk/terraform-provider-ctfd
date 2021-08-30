package provider

import (
	"context"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
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

	team := new(api.NewTeam)
	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}

	currentTeam, err := client.GetTeam(uint(intId))
	if err != nil {
		return diag.FromErr(err)
	}
	team.Password = d.Get("password").(string)
	if d.HasChange("name") {
		team.Name = d.Get("name").(string)
	} else {
		team.Name = currentTeam.Name
	}
	if d.HasChange("email") {
		team.Email = d.Get("email").(string)
	} else {
		team.Email = currentTeam.Email
	}
	if d.HasChange("website") {
		team.Website = d.Get("website").(string)
	} else {
		team.Website = currentTeam.Website
	}
	if d.HasChange("affiliation") {
		team.Affiliation = d.Get("affiliation").(string)
	} else {
		team.Affiliation = currentTeam.Affiliation
	}
	if d.HasChange("country") {
		team.Country = d.Get("country").(string)
	} else {
		team.Country = currentTeam.Country
	}
	if d.HasChange("hidden") {
		team.Hidden = d.Get("hidden").(bool)
	} else {
		team.Hidden = currentTeam.Hidden
	}
	if d.HasChange("banned") {
		team.Banned = d.Get("banned").(bool)
	} else {
		team.Banned = currentTeam.Banned
	}

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
			"password": &schema.Schema{
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
			},
		},
	}
}
