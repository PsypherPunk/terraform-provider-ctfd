package provider

import (
	"context"

	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandCtfdSetupEmailConfig(l []interface{}) *api.EmailConfig {
	if len(l) == 0 {
		return nil
	}

	m := l[0].(map[string]interface{})

	emailConfig := &api.EmailConfig{
		Username:    m["username"].(string),
		Password:    m["password"].(string),
		FromAddress: m["from_address"].(string),
		Server:      m["server"].(string),
		Port:        m["port"].(int),
		UseAuth:     m["use_auth"].(bool),
		UseTls:      m["use_tls"].(bool),
		UseSsl:      m["use_ssl"].(bool),
	}

	return emailConfig
}

func resourceCtfdSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	setup := api.CtfdSetup{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		AdminEmail:        d.Get("admin_email").(string),
		ConfigurationPath: d.Get("configuration_path").(string),
	}

	if v, ok := d.GetOk("email"); ok {
		setup.Email = expandCtfdSetupEmailConfig(v.([]interface{}))
	}

	err := client.CreateCtfdSetup(setup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(setup.Name)

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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"from_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"use_auth": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"use_tls": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"use_ssl": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}
