package provider

import (
	"context"
	"log"

	sharedmodels "github.com/hashicorp/cloud-sdk-go/clients/cloud-shared/v1/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-hcp/internal/clients"
)

func dataSourceHvn() *schema.Resource {
	return &schema.Resource{
		Description: "The HVN data source provides information about an existing HashiCorp Virtual Network.",
		ReadContext: dataSourceHvnRead,
		Timeouts: &schema.ResourceTimeout{
			Default: &hvnDefaultTimeout,
		},
		Schema: map[string]*schema.Schema{
			// Required inputs
			"hvn_id": {
				Description:      "The ID of the HashiCorp Virtual Network.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateSlugID,
			},
			"cloud_provider": {
				Description:      "The provider where the HVN is located.",
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateStringInSlice(hvnResourceCloudProviders, true),
			},
			"region": {
				Description: "The region where the HVN is located.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Computed outputs
			"cidr_block": {
				Description: "The CIDR range of the HVN.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization_id": {
				Description: "The ID of the HCP organization where the HVN is located.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"project_id": {
				Description: "The ID of the HCP project where the HVN is located.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "The time that the HVN was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// dataSourceHvnRead is the func to implement reading of an
// HashiCorp Virtual Network (HVN)
func dataSourceHvnRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*clients.Client)

	hvnID := d.Get("hvn_id").(string)

	loc := &sharedmodels.HashicorpCloudLocationLocation{
		OrganizationID: client.Config.OrganizationID,
		ProjectID:      client.Config.ProjectID,
		Region: &sharedmodels.HashicorpCloudLocationRegion{
			Provider: d.Get("cloud_provider").(string),
			Region:   d.Get("region").(string),
		},
	}

	// Check for an existing HVN
	log.Printf("[INFO] Reading HVN (%s) [project_id=%s, organization_id=%s]", hvnID, loc.ProjectID, loc.OrganizationID)
	hvn, err := clients.GetHvnByID(ctx, client, loc, hvnID)
	if err != nil {
		return diag.FromErr(err)
	}

	link := newLink(loc, HvnResourceType, hvnID)
	url, err := linkURL(link)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(url)

	if err := setHvnResourceData(d, hvn); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
