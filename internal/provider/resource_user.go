package provider

import (
	"context"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	user := api.NewUser{
		Name:        d.Get("name").(string),
		Email:       d.Get("email").(string),
		Password:    d.Get("password").(string),
		Website:     d.Get("website").(string),
		Affiliation: d.Get("affiliation").(string),
		Country:     d.Get("country").(string),
		Type:        d.Get("type").(string),
		Verified:    d.Get("verified").(bool),
		Hidden:      d.Get("hidden").(bool),
		Banned:      d.Get("banned").(bool),
	}

	newUser, err := client.CreateUser(user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(newUser.Id)))
	if err != nil {
		return nil
	}
	err = d.Set("fields", newUser.Fields)
	if err != nil {
		return nil
	}

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Get("id").(string)

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	user, err := client.GetUser(uint(intId))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(user.Id)))

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Id()

	intId, err := strconv.Atoi(id)
	if err != nil {
		return diag.FromErr(err)
	}

	user := new(api.NewUser)
	user.Password = d.Get("password").(string)
	user.Name = d.Get("name").(string)
	user.Email = d.Get("email").(string)
	user.Website = d.Get("website").(string)
	user.Affiliation = d.Get("affiliation").(string)
	user.Country = d.Get("country").(string)
	user.Verified = d.Get("verified").(bool)
	user.Hidden = d.Get("hidden").(bool)
	user.Banned = d.Get("banned").(bool)
	user.Type = d.Get("type").(string)

	updatedUser, err := client.UpdateUser(uint(intId), *user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(updatedUser.Id)))

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	id := d.Id()

	var diags diag.Diagnostics

	intId, err := strconv.Atoi(id)
	if err != nil {
		return nil
	}
	err = client.DeleteUser(uint(intId))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Get details of a user.",
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
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
			"verified": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
}
