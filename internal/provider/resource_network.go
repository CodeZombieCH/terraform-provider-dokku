package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/melbahja/goph"
)

func resourceNetwork() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages a docker network resource.",
		CreateContext: resourceNetworkCreate,
		ReadContext:   resourceNetworkRead,
		DeleteContext: resourceNetworkDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the network.",
				ForceNew:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNetworkCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sshClient := m.(*goph.Client)

	var diags diag.Diagnostics

	network := NewDokkuNetworkFromResourceData(d)
	err := dokkuNetworkCreate(network, sshClient)
	if err != nil {
		return diag.FromErr(err)
	}

	err = dokkuNetworkRead(network, sshClient)
	if err != nil {
		return diag.FromErr(err)
	}

	network.setOnResourceData(d)

	return diags
}

func resourceNetworkRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sshClient := m.(*goph.Client)

	var diags diag.Diagnostics

	var serviceName string
	if d.Id() == "" {
		serviceName = d.Get("name").(string)
	} else {
		serviceName = d.Id()
	}

	network := NewDokkuNetwork(serviceName)
	err := dokkuNetworkRead(network, sshClient)
	if err != nil {
		return diag.FromErr(err)
	}

	network.setOnResourceData(d)

	return diags
}

func resourceNetworkUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sshClient := m.(*goph.Client)

	var diags diag.Diagnostics

	network := NewDokkuNetworkFromResourceData(d)
	err := dokkuNetworkUpdate(network, d, sshClient)

	if err != nil {
		return diag.FromErr(err)
	}

	network.setOnResourceData(d)

	return diags
}

func resourceNetworkDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	sshClient := m.(*goph.Client)

	var diags diag.Diagnostics

	network := NewDokkuNetworkFromResourceData(d)
	err := dokkuNetworkDestroy(network, sshClient)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
