package catalog_test

import (
	"fmt"
	"testing"

	"github.com/databricks/terraform-provider-databricks/internal/acceptance"
)

func TestUcAccCatalog(t *testing.T) {
	acceptance.LoadUcwsEnv(t)
	acceptance.UnityWorkspaceLevel(t, acceptance.Step{
		Template: fmt.Sprintf(`
		resource "databricks_catalog" "sandbox" {
			name         = "sandbox{var.RANDOM}"
			comment      = "this catalog is managed by terraform"
			properties = {
				purpose = "testing"
			}
			%s
		}`, getPredictiveOptimizationSetting(t, true)),
	})
}

func TestUcAccCatalogIsolated(t *testing.T) {
	acceptance.UnityWorkspaceLevel(t, acceptance.Step{
		Template: `
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
		}`,
	}, acceptance.Step{
		Template: `
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			isolation_mode = "ISOLATED"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
		}`,
	}, acceptance.Step{
		Template: `
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			isolation_mode = "OPEN"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
		}`,
	})
}

func TestUcAccCatalogUpdate(t *testing.T) {
	acceptance.LoadUcwsEnv(t)
	acceptance.UnityWorkspaceLevel(t, acceptance.Step{
		Template: fmt.Sprintf(`
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
			%s
		}`, getPredictiveOptimizationSetting(t, true)),
	}, acceptance.Step{
		Template: fmt.Sprintf(`
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
			%s
			owner = "account users"
		}`, getPredictiveOptimizationSetting(t, true)),
	}, acceptance.Step{
		Template: fmt.Sprintf(`
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			comment        = "this catalog is managed by terraform"
			properties     = {
				purpose = "testing"
			}
			%s
			owner = "{env.TEST_DATA_ENG_GROUP}"
		}`, getPredictiveOptimizationSetting(t, true)),
	}, acceptance.Step{
		Template: fmt.Sprintf(`
		resource "databricks_catalog" "sandbox" {
			name           = "sandbox{var.STICKY_RANDOM}"
			comment        = "this catalog is managed by terraform - updated comment"
			properties     = {
				purpose = "testing"
			}
			%s
			owner = "{env.TEST_METASTORE_ADMIN_GROUP_NAME}"
		}`, getPredictiveOptimizationSetting(t, false)),
	})
}
