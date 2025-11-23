package provider

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/blang/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/melbahja/goph"
)

type DokkuNetwork struct {
	Name   string
	Driver string
}

func NewDokkuNetwork(name string) *DokkuNetwork {
	return &DokkuNetwork{
		Name: name,
	}
}

func NewDokkuNetworkFromResourceData(d *schema.ResourceData) *DokkuNetwork {
	return &DokkuNetwork{
		Name: d.Get("name").(string),
	}
}

func (s *DokkuNetwork) setOnResourceData(d *schema.ResourceData) {
	d.SetId(s.Name)
	d.Set("name", s.Name)
}

type dokkuNetworkInfoResponse struct {
	Name   string                         `json:"Name"`
	Driver string                         `json:"Driver"`
	Labels dokkuNetworkInfoLabelsResponse `json:"Labels"`
}

type dokkuNetworkInfoLabelsResponse struct {
	DokkuNetworkName string `json:"com.dokku.network-name"`
}

func dokkuNetworkExists(name string, client *goph.Client) bool {
	res := run(client, fmt.Sprintf("network:exists %s", name))
	if res.status == 0 {
		return true
	} else {
		return false
	}
}

func dokkuNetworkGetInfo(network *DokkuNetwork, client *goph.Client) error {
	res := run(client, fmt.Sprintf("network:info %s --format json", network.Name))

	if res.err != nil {
		return res.err
	}

	var response dokkuNetworkInfoResponse
	if err := json.Unmarshal([]byte(res.stdout), &response); err != nil {
		return err
	}

	network.Driver = response.Driver
	return nil
}

func dokkuNetworkRetrieve(name string, client *goph.Client) (*DokkuNetwork, error) {
	network := &DokkuNetwork{Name: name}
	err := dokkuNetworkRead(network, client)
	return network, err
}

func dokkuNetworkRead(network *DokkuNetwork, client *goph.Client) error {
	// The `network:info` command was introduced in v0.35.3. If a lower version is used,
	// fall back to `network:exists`
	// See https://dokku.com/docs/networking/network/#checking-network-info
	infoCommandAvailableRange := ">= 0.35.3"
	infoCommandAvailable, err := semver.ParseRange(infoCommandAvailableRange)
	if err != nil {
		return fmt.Errorf("failed to parse semver to check `network:info` command availability")
	}

	if infoCommandAvailable(DOKKU_VERSION) {
		return dokkuNetworkGetInfo(network, client)
	} else {
		if exists := dokkuNetworkExists(network.Name, client); !exists {
			return fmt.Errorf("network %s does not exist", network.Name)
		}
		// If it exists, nothing to do here, as name is the only attibute of interest, and it's already set
	}

	return nil
}

func dokkuNetworkCreate(network *DokkuNetwork, client *goph.Client) error {
	res := run(client, fmt.Sprintf("network:create %s", network.Name))

	log.Printf("[DEBUG] network:create %v\n", res.stdout)

	if res.err != nil {
		return res.err
	}

	return nil
}

func dokkuNetworkUpdate(network *DokkuNetwork, d *schema.ResourceData, client *goph.Client) error {
	return fmt.Errorf("not supported")
}

func dokkuNetworkDestroy(network *DokkuNetwork, client *goph.Client) error {
	res := run(client, fmt.Sprintf("network:destroy %s --force", network.Name))

	log.Printf("[DEBUG] network:destroy %v\n", res.stdout)

	if res.err != nil {
		return res.err
	}

	return nil
}
