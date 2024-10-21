package repos

import (
	"context"
	"fmt"
	"github.com/databricks/databricks-sdk-go/service/workspace"
	"github.com/databricks/terraform-provider-databricks/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/url"
	"regexp"
	"strings"
)

type gitFolderStruct struct {
	workspace.RepoInfo
}

var gitFolderAliasMap = map[string]string{
	"provider": "git_provider",
}

func (gitFolderStruct) Aliases() map[string]map[string]string {
	return map[string]map[string]string{
		"workspace.gitFolderStruct": gitFolderAliasMap,
	}
}

type gitFolderCreateStruct struct {
	workspace.CreateRepoRequest
}

func (gitFolderCreateStruct) Aliases() map[string]map[string]string {
	return map[string]map[string]string{
		"workspace.gitFolderCreateStruct": gitFolderAliasMap,
	}
}

func (gitFolderCreateStruct) CustomizeSchema(s *common.CustomizableSchema) *common.CustomizableSchema {
	return s
}

type gitFolderUpdateStruct struct {
	workspace.UpdateRepoRequest
}

func (gitFolderUpdateStruct) Aliases() map[string]map[string]string {
	return map[string]map[string]string{
		"workspace.gitFolderUpdateStruct": gitFolderAliasMap,
	}
}

func (gitFolderUpdateStruct) CustomizeSchema(s *common.CustomizableSchema) *common.CustomizableSchema {
	return s
}

var (
	gitProvidersMap = map[string]string{
		"github.com":    "gitHub",
		"dev.azure.com": "azureDevOpsServices",
		"gitlab.com":    "gitLab",
		"bitbucket.org": "bitbucketCloud",
	}
	awsCodeCommitRegex = regexp.MustCompile(`^git-codecommit\.[^.]+\.amazonaws\.com$`)
)

func ResourceGitFolder() common.Resource {
	s := common.StructToSchema(gitFolderStruct{}, func(m map[string]*schema.Schema) map[string]*schema.Schema {
		common.CustomizeSchemaPath(m).RemoveField("id")
		common.CustomizeSchemaPath(m).AddNewField("id", &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		})

		for _, p := range []string{"git_provider", "path", "branch", "head_commit_id"} {
			common.CustomizeSchemaPath(m, p).SetComputed()
		}
		for _, p := range []string{"url", "git_provider", "sparse_checkout", "path"} {
			common.CustomizeSchemaPath(m, p).SetForceNew()
		}
		return m
	})
	return common.Resource{
		Create: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			var create gitFolderCreateStruct
			if create.Provider == "" { // trying to infer Git Provider from the URL
				create.Provider = GetGitProviderFromUrl(create.Url)
			}
			if create.Provider == "" {
				return fmt.Errorf("git_provider isn't specified and we can't detect provider from URL")
			}
			common.DataToStructPointer(d, s, &create)
			ws, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			repo, err := ws.Repos.Create(ctx, workspace.CreateRepoRequest{
				Path:           create.Path,
				Provider:       create.Provider,
				SparseCheckout: create.SparseCheckout,
				Url:            create.Url,
			})
			if err != nil {
				return err
			}
			//d.Set("id", repo.Id)
			common.StructToData(repo, s, d)
			return nil
		},
		Read: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			ws, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			repo, err := ws.Repos.GetByRepoId(ctx, 0)
			if err != nil {
				return err
			}
			return common.StructToData(repo, s, d)
		},
		Update: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			var update gitFolderUpdateStruct
			common.DataToStructPointer(d, s, &update)
			ws, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			repo := ws.Repos.Update(ctx, workspace.UpdateRepoRequest{
				Branch:         update.Branch,
				RepoId:         update.RepoId,
				SparseCheckout: update.SparseCheckout,
				Tag:            update.Tag,
			})
			return common.StructToData(repo, s, d)
		},
		Delete: func(ctx context.Context, d *schema.ResourceData, c *common.DatabricksClient) error {
			ws, err := c.WorkspaceClient()
			if err != nil {
				return err
			}
			return ws.Repos.DeleteByRepoId(ctx, 0)
		},
		Schema: s,
	}
}

func GetGitProviderFromUrl(uri string) string {
	provider := ""
	u, err := url.Parse(uri)
	if err == nil {
		lhost := strings.ToLower(u.Host)
		provider = gitProvidersMap[lhost]
		if provider == "" && awsCodeCommitRegex.FindStringSubmatch(lhost) != nil {
			provider = "awsCodeCommit"
		}
	}
	return provider
}
