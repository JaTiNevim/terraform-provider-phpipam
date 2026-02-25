package phpipam

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type PingResponse struct {
	ScanType   string `json:"scan_type"`
	ExitCode   int    `json:"exit_code"`
	ResultCode string `json:"result_code"`
	Message    string `json:"message"`
}

func dataSourcePHPIPAMPing() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePHPIPAMPingRead,
		Schema: map[string]*schema.Schema{
			"address_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the IP address to ping.",
			},
			"ping_result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The result of the ping operation.",
			},
		},
	}
}

func dataSourcePHPIPAMPingRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*ProviderPHPIPAMClient).addressesController

	// Verifying that the client is properly initialized
	if c == nil {
		return fmt.Errorf("addressesController is nil")
	}

	var pingResult PingResponse
	addressID := d.Get("address_id").(int)

	err := c.SendRequest("GET", fmt.Sprintf("/addresses/%d/ping/", addressID), &struct{}{}, &pingResult)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid Id") {
			d.SetId("")
			return nil
		}
		return err
	}

	// Setting ping result to Terraform
	if err := d.Set("ping_result", pingResult.ResultCode); err != nil {
		return err
	}

	// Setting the Adress ID
	d.SetId(strconv.Itoa(addressID))

	return nil
}
