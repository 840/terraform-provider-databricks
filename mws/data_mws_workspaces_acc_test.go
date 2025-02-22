package mws_test

import (
	"fmt"

	"github.com/databricks/terraform-provider-databricks/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"testing"
)

func TestAccDataSourceMwsWorkspaces(t *testing.T) {
	acceptance.AccountLevel(t, acceptance.Step{
		Template: `
		data "databricks_mws_workspaces" "this" {
		}`,
		Check: func(s *terraform.State) error {
			r, ok := s.RootModule().Resources["data.databricks_mws_workspaces.this"]
			if !ok {
				return fmt.Errorf("data not found in state")
			}
			ids := r.Primary.Attributes["ids.%"]
			if ids == "" {
				return fmt.Errorf("ids is empty: %v", r.Primary.Attributes)
			}
			return nil
		},
	})
}
