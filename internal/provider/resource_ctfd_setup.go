package provider

import (
	"context"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCtfdSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	setup := api.CtfdSetup{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AdminEmail:        d.Get("admin_email").(string),
		ConfigurationPath: d.Get("configuration_path").(string),
	}

	err := client.CreateCtfdSetup(setup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(setup.Name)
	if err != nil {
		return nil
	}

	return diags
}

func resourceCtfdSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	setup, err := client.GetCtfdSetup()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(setup.Name)

	return diags
}

func resourceCtfdSetupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	err := client.DeleteCtfdSetup()
	if err != nil {
		return diag.FromErr(err)
	}

	setup := new(api.CtfdSetup)
	setup.Name = d.Get("name").(string)
	setup.Description = d.Get("description").(string)
	setup.ConfigurationPath = d.Get("configuration_path").(string)

	err = client.CreateCtfdSetup(*setup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(setup.Name)

	return diags
}

func resourceCtfdSetupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	err := client.DeleteCtfdSetup()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceCtfdSetup() *schema.Resource {
	return &schema.Resource{
		Description:   "Initial setup for a CTFd instance.",
		CreateContext: resourceCtfdSetupCreate,
		ReadContext:   resourceCtfdSetupRead,
		UpdateContext: resourceCtfdSetupUpdate,
		DeleteContext: resourceCtfdSetupDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"admin_email": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_path": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
