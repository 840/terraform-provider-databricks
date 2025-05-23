package vectorsearch

import (
	"fmt"
	"testing"
	"time"

	"github.com/databricks/databricks-sdk-go/experimental/mocks"
	"github.com/databricks/databricks-sdk-go/qa/poll"
	"github.com/databricks/terraform-provider-databricks/qa"

	"github.com/databricks/databricks-sdk-go/service/vectorsearch"

	"github.com/stretchr/testify/mock"
)

func TestVectorSearchEndpointCornerCases(t *testing.T) {
	qa.ResourceCornerCases(t, ResourceVectorSearchEndpoint())
}

func TestVectorSearchEndpointCreate(t *testing.T) {
	ei := &vectorsearch.EndpointInfo{
		Name:           "abc",
		EndpointStatus: &vectorsearch.EndpointStatus{State: "ONLINE"},
		Id:             "1234-5678",
	}
	qa.ResourceFixture{
		MockWorkspaceClientFunc: func(w *mocks.MockWorkspaceClient) {
			e := w.GetMockVectorSearchEndpointsAPI().EXPECT()
			e.CreateEndpoint(mock.Anything, vectorsearch.CreateEndpoint{
				Name:         "abc",
				EndpointType: "STANDARD",
			}).Return(&vectorsearch.WaitGetEndpointVectorSearchEndpointOnline[vectorsearch.EndpointInfo]{Poll: poll.Simple(*ei)}, nil)
			e.GetEndpointByEndpointName(mock.Anything, "abc").Return(ei, nil)
		},
		Resource: ResourceVectorSearchEndpoint(),
		HCL: `
		name          = "abc"
		endpoint_type = "STANDARD"
		`,
		Create: true,
	}.ApplyAndExpectData(t, map[string]any{
		"id":          "abc",
		"endpoint_id": "1234-5678",
	})
}

func TestVectorSearchEndpointRead(t *testing.T) {
	ei := &vectorsearch.EndpointInfo{
		Name:           "abc",
		EndpointStatus: &vectorsearch.EndpointStatus{State: "ONLINE"},
		Id:             "1234-5678",
	}
	qa.ResourceFixture{
		MockWorkspaceClientFunc: func(w *mocks.MockWorkspaceClient) {
			e := w.GetMockVectorSearchEndpointsAPI().EXPECT()
			e.GetEndpointByEndpointName(mock.Anything, "abc").Return(ei, nil)
		},
		Resource: ResourceVectorSearchEndpoint(),
		ID:       "abc",
		HCL: `
		name          = "abc"
		endpoint_type = "STANDARD"
		`,
		Read: true,
	}.ApplyAndExpectData(t, map[string]any{
		"id":          "abc",
		"endpoint_id": "1234-5678",
	})
}

func TestVectorSearchEndpointDelete(t *testing.T) {
	qa.ResourceFixture{
		MockWorkspaceClientFunc: func(w *mocks.MockWorkspaceClient) {
			w.GetMockVectorSearchEndpointsAPI().EXPECT().DeleteEndpointByEndpointName(mock.Anything, "abc").Return(nil)
		},
		Resource: ResourceVectorSearchEndpoint(),
		Delete:   true,
		ID:       "abc",
	}.ApplyAndExpectData(t, map[string]any{
		"id": "abc",
	})
}

func TestVectorSearchEndpointCreateTimeoutError(t *testing.T) {
	qa.ResourceFixture{
		MockWorkspaceClientFunc: func(w *mocks.MockWorkspaceClient) {
			e := w.GetMockVectorSearchEndpointsAPI().EXPECT()
			e.CreateEndpoint(mock.Anything, vectorsearch.CreateEndpoint{
				Name:         "abc",
				EndpointType: "STANDARD",
			}).Return(&vectorsearch.WaitGetEndpointVectorSearchEndpointOnline[vectorsearch.EndpointInfo]{
				Poll: func(_ time.Duration, _ func(*vectorsearch.EndpointInfo)) (*vectorsearch.EndpointInfo, error) {
					return nil, fmt.Errorf("timeout")
				},
			}, nil)
			e.DeleteEndpointByEndpointName(mock.Anything, "abc").Return(nil)
		},
		Resource: ResourceVectorSearchEndpoint(),
		HCL: `
		name          = "abc"
		endpoint_type = "STANDARD"
		`,
		Create: true,
	}.ExpectError(t, "timeout")

}

func TestVectorSearchEndpointUpdateBudgetPolicy(t *testing.T) {
	qa.ResourceFixture{
		MockWorkspaceClientFunc: func(w *mocks.MockWorkspaceClient) {
			e := w.GetMockVectorSearchEndpointsAPI().EXPECT()
			e.UpdateEndpointBudgetPolicy(mock.Anything, vectorsearch.PatchEndpointBudgetPolicyRequest{
				EndpointName:   "abc",
				BudgetPolicyId: "budget-123",
			}).Return(&vectorsearch.PatchEndpointBudgetPolicyResponse{
				EffectiveBudgetPolicyId: "budget-123",
			}, nil)
			e.GetEndpointByEndpointName(mock.Anything, "abc").Return(&vectorsearch.EndpointInfo{
				Id:                      "1234-5678",
				EffectiveBudgetPolicyId: "budget-123",
				EndpointStatus:          &vectorsearch.EndpointStatus{State: "ONLINE"},
				EndpointType:            "STANDARD",
				Name:                    "abc",
			}, nil)
		},
		Resource: ResourceVectorSearchEndpoint(),
		Update:   true,
		ID:       "abc",
		InstanceState: map[string]string{
			"endpoint_id":      "1234-5678",
			"name":             "abc",
			"endpoint_type":    "STANDARD",
			"budget_policy_id": "budget-456",
		},
		HCL: `
		name             = "abc"
		endpoint_type    = "STANDARD"
		budget_policy_id = "budget-123"
		`,
	}.ApplyAndExpectData(t, map[string]any{
		"id":               "abc",
		"budget_policy_id": "budget-123",
	})
}
