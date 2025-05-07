package atlantis

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysaidit/atlantis-drift-detection/internal/testhelper"
)

func makeTestClient(t *testing.T) *Client {
	c := Client{
		AtlantisHostname: testhelper.EnvOrSkip(t, "ATLANTIS_HOST"),
		Token:            testhelper.EnvOrSkip(t, "ATLANTIS_TOKEN"),
		HTTPClient:       http.DefaultClient,
	}
	return &c
}

func loadPlanSummaryOk(t *testing.T) *PlanSummaryRequest {
	return loadPlanSummaryRequest(t, "PLAN_SUMMARY_OK")
}

func loadPlanSummaryRequest(t *testing.T, name string) *PlanSummaryRequest {
	body := testhelper.EnvOrSkip(t, name)
	var ret PlanSummaryRequest
	require.NoError(t, json.Unmarshal([]byte(body), &ret))
	return &ret
}

func loadPlanSummaryLock(t *testing.T) *PlanSummaryRequest {
	return loadPlanSummaryRequest(t, "PLAN_SUMMARY_LOCK")
}

func loadPlanSummaryChanges(t *testing.T) *PlanSummaryRequest {
	return loadPlanSummaryRequest(t, "PLAN_SUMMARY_CHANGES")
}

func TestClient_PlanSummaryLock(t *testing.T) {
	testhelper.ReadEnvFile(t, "../../")
	c := makeTestClient(t)
	req := loadPlanSummaryLock(t)
	ok, err := c.PlanSummary(context.Background(), req)
	require.NoError(t, err)
	require.True(t, ok.IsLocked())
}

func TestClient_PlanSummaryOk(t *testing.T) {
	testhelper.ReadEnvFile(t, "../../")
	c := makeTestClient(t)
	req := loadPlanSummaryOk(t)
	ok, err := c.PlanSummary(context.Background(), req)
	require.NoError(t, err)
	require.False(t, ok.HasChanges())
}

func TestClient_PlanSummaryChanges(t *testing.T) {
	testhelper.ReadEnvFile(t, "../../")
	c := makeTestClient(t)
	req := loadPlanSummaryChanges(t)
	ok, err := c.PlanSummary(context.Background(), req)
	require.NoError(t, err)
	require.True(t, ok.HasChanges())
}
