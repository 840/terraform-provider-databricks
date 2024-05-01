package scim

import (
	"context"
	"fmt"
	"github.com/databricks/databricks-sdk-go"
	"sort"

	"github.com/databricks/databricks-sdk-go/service/iam"
	"github.com/databricks/terraform-provider-databricks/common"
)

func DataSourceUsers() common.Resource {
	type usersData struct {
		DisplayNameContains string   `json:"display_name_contains,omitempty" tf:"computed"`
		ID                  []string `json:"ids,omitempty" tf:"computed,slice_set"`
	}
	return common.AccountData(func(ctx context.Context, d *usersData, c *databricks.AccountClient) error {
		response := d
		usersList, err := c.Users.ListAll(ctx, iam.ListAccountUsersRequest{
			Filter: fmt.Sprintf(`displayName co "%s"`, response.DisplayNameContains),
		})
		if err != nil {
			return err
		}
		if len(usersList) == 0 {
			return fmt.Errorf("cannot find users with display name containing %s", response.DisplayNameContains)
		}
		for _, user := range usersList {
			response.ID = append(response.ID, user.Id)
		}
		sort.Strings(response.ID)
		return nil
	})
}
