package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				"username": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Admin. username",
				},
				"password": &schema.Schema{
					Type:        schema.TypeString,
					Required:    true,
					Description: "Admin. password",
				},
				"url": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Base URL of CTFd instance",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"ctfd_ctfd_instance": dataSourceCtfd(),
				"ctfd_challenges":    dataSourceChallenges(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"ctfd_resource": resourceCtfd(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	httpClient    *http.Client
	httpUserAgent string
	ctfdUsername  string
	ctfdPassword  string
	ctfdUrl       string
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		ctfdUrl := d.Get("url").(string)
		ctfdUsername := d.Get("username").(string)
		ctfdPassword := d.Get("password").(string)
		userAgent := p.UserAgent("terraform-provider-ctfd", version)

		httpClient := http.DefaultClient

		return &apiClient{
			httpClient:    httpClient,
			httpUserAgent: userAgent,
			ctfdUsername:  ctfdUsername,
			ctfdPassword:  ctfdPassword,
			ctfdUrl:       ctfdUrl,
		}, nil
	}
}
