package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/PsypherPunk/terraform-provider-ctfd/internal/api"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"username": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Admin. username",
				},
				"password": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Admin. password",
				},
				"url": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Base URL of CTFd instance",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"ctfd_challenges": dataSourceChallenges(),
				"ctfd_teams":      dataSourceTeams(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"ctfd_setup":                resourceCtfdSetup(),
				"ctfd_team":                 resourceTeam(),
				"ctfd_user":                 resourceUser(),
				"ctfd_user_team_membership": resourceUserTeamMembership(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		ctfdUrl := d.Get("url").(string)
		ctfdUsername := d.Get("username").(string)
		ctfdPassword := d.Get("password").(string)
		userAgent := p.UserAgent("terraform-provider-ctfd", version)

		client, err := api.NewClient(&ctfdUrl, &ctfdUsername, &ctfdPassword, &userAgent)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		err = client.CheckSetup()
		if err == nil {
			err = client.SignIn()
			if err != nil {
				return nil, diag.FromErr(err)
			}

			token, err := client.GetOrCreateToken()
			if err != nil {
				return nil, diag.FromErr(err)
			}
			client.Auth.Token = token.Value
		}

		return client, nil
	}
}
