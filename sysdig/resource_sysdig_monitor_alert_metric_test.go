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

func TestAccAlertMetric(t *testing.T) {
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
				Config: alertMetricWithName(rText()),
			},
			{
				Config: alertMetricWithoutScopeWithName(rText()),
			},
			{
				Config: alertMetricWithNotificationChannel(rText()),
			},
		},
	})
}

func alertMetricWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	scope = "kubernetes.cluster.name in (\"pulsar\")"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

func alertMetricWithoutScopeWithName(name string) string {
	return fmt.Sprintf(`
resource "sysdig_monitor_alert_metric" "sample2" {
	name = "TERRAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 3

	metric = "avg(avg(cpu.used.percent)) > 50"
	
	trigger_after_minutes = 10

	enabled = false

	multiple_alerts_by = ["kubernetes.deployment.name"]

	capture {
		filename = "TERRAFORM_TEST.scap"
		duration = 15
	}
}
`, name, name)
}

// Reported by @logdnalf at https://github.com/draios/terraform-provider-sysdig/issues/24
func alertMetricWithNotificationChannel(name string) string {
	return fmt.Sprintf(`
resource "sysdig_secure_notification_channel" "sample-pagerduty" {
	name = "Example Channel %s - Pagerduty"
	enabled = true
	type = "PAGER_DUTY"
	account = "account"
	service_key = "XXXXXXXXXX"
	service_name = "sysdig"
	notify_when_ok = true
	notify_when_resolved = true
}

resource "sysdig_monitor_alert_metric" "sample3" {
	enabled = true
	name = "TERAFORM TEST - METRIC %s"
	description = "TERRAFORM TEST - METRIC %s"
	severity = 6
	metric = "sum(min(cpu.used.percent)) > 100000"
	scope = "kubernetes.cluster.name in (\"foo\")"
	trigger_after_minutes = 20
	notification_channels = [
	sysdig_secure_notification_channel.sample-pagerduty.id
	]
	multiple_alerts_by = [
	"kubernetes.cluster.name"
	]
}`, name, name, name)
}
