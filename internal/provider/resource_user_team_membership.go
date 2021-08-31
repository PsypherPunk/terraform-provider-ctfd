package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

func resourceUserTeamMembershipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	user_id := d.Get("user_id").(int)
	team_id := d.Get("team_id").(int)

	newUserTeamMembership, err := client.CreateUserTeamMembership(uint(team_id), uint(user_id))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d.%d", newUserTeamMembership.TeamId, newUserTeamMembership.UserId))

	return diags
}

func resourceUserTeamMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Id()
	ids := strings.SplitN(id, ".", 2)
	teamId := ids[0]
	userId := ids[1]

	teamIntId, err := strconv.Atoi(teamId)
	userIntId, err := strconv.Atoi(userId)

	members, err := client.GetTeamMemberships(uint(teamIntId))
	if err != nil {
		return diag.FromErr(err)
	}

	if !api.Contains(*members, uint(userIntId)) {
		return diag.FromErr(errors.New(fmt.Sprintf("member %d not found in team %d", userIntId, teamIntId)))
	}

	d.SetId(id)

	return diags
}

func resourceUserTeamMembershipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Id()
	ids := strings.SplitN(id, ".", 2)
	teamId := ids[0]
	userId := ids[1]

	teamIntId, err := strconv.Atoi(teamId)
	userIntId, err := strconv.Atoi(userId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteUserTeamMembership(uint(teamIntId), uint(userIntId))
	if err != nil {
		return diag.FromErr(err)
	}

	teamUintId := d.Get("user_id").(int)
	userUintId := d.Get("team_id").(int)

	newUserTeamMembership, err := client.CreateUserTeamMembership(uint(teamUintId), uint(userUintId))
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d.%d", newUserTeamMembership.TeamId, newUserTeamMembership.UserId))

	return diags
}

func resourceUserTeamMembershipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*api.Client)

	var diags diag.Diagnostics

	id := d.Id()
	ids := strings.SplitN(id, ".", 2)
	teamId := ids[0]
	userId := ids[1]

	teamIntId, err := strconv.Atoi(teamId)
	userIntId, err := strconv.Atoi(userId)
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteUserTeamMembership(uint(teamIntId), uint(userIntId))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserTeamMembership() *schema.Resource {
	return &schema.Resource{
		Description:   "Get details of a User/Team Membership.",
		CreateContext: resourceUserTeamMembershipCreate,
		ReadContext:   resourceUserTeamMembershipRead,
		UpdateContext: resourceUserTeamMembershipUpdate,
		DeleteContext: resourceUserTeamMembershipDelete,
		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"team_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}
