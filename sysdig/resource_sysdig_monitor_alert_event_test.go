package sysdig_test

import (
	"fmt"
	"github.com/draios/terraform-provider-sysdig/sysdig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"os"
	"testing"
)

func TestAccAlertEvent(t *testing.T) {
	rText := func() string { return acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum) }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			if v := os.Getenv("SYSDIG_MONITOR_API_TOKEN"); v == "" {
				t.Fatal("SYSDIG_MONITOR_API_TOKEN must be set for acceptance tests")
			}
		},
		Providers: map[string]terraform.ResourceProvider{
			"sysdig": sysdig.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: alertEventWithName(rText()),
			},
		},
	})
}

func alertEventWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_event" "sample" {
	name = "TERRAFORM TEST - EVENT %s"
	description = "TERRAFORM TEST - EVENT %s"
	severity = 4

	event_name = "deployment"
	source = "kubernetes"
	event_rel = ">"
	event_count = 2

	multiple_alerts_by = ["kubernetes.deployment.name"]
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false
}
`, name, name)
}
