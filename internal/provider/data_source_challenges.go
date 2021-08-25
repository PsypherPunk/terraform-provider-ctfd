package provider

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type APIResult struct {
	Success bool             `json:"success"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

func dataSourceChallengesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)

	var diags diag.Diagnostics

	//url := strings.TrimRight(client.ctfdUrl, "/")
	jar, err := cookiejar.New(nil)
	if err != nil {
		return diag.FromErr(err)
	}
	client.httpClient.Jar = jar

	nonceRegex := regexp.MustCompile("'csrfNonce': \"([a-z0-9]+)\",")

	// GET nonce
	resp, err := client.httpClient.Get("http://0.0.0.0:8000/setup")
	if err != nil {
		return diag.FromErr(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	parts := nonceRegex.FindSubmatch(body)
	nonce := parts[1]

	// Initial auth.
	form := url.Values{}
	form.Set("nonce", string(nonce))
	form.Set("name", client.ctfdUsername)
	form.Set("password", client.ctfdPassword)
	resp, err = client.httpClient.PostForm("http://0.0.0.0:8000/login", form)
	if err != nil {
		return diag.FromErr(err)
	}

	// GET nonce
	resp, err = client.httpClient.Get("http://0.0.0.0:8000/")
	if err != nil {
		return diag.FromErr(err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	parts = nonceRegex.FindSubmatch(body)
	nonce = parts[1]

	// GET Challenges
	req, err := http.NewRequest("GET", "http://0.0.0.0:8000/api/v1/challenges", nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("CSRF-Token", string(nonce))

	resp, err = client.httpClient.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}

	result := new(APIResult)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return diag.FromErr(err)
	}

	challenges := make([]map[string]interface{}, 0)
	err = json.Unmarshal(*result.Data, &challenges)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("challenges", challenges); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceChallenges() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of the current challenges.",
		ReadContext: dataSourceChallengesRead,
		Schema: map[string]*schema.Schema{
			"challenges": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solves": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"solved_by_me": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"category": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"template": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"script": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
