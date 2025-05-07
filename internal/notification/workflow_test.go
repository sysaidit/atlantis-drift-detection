package notification

import (
	"context"
	"testing"

	"github.com/cresta/gogithub"
	"github.com/sysaidit/atlantis-drift-detection/internal/testhelper"
	"go.uber.org/zap/zaptest"
)

func TestNewWorkflow(t *testing.T) {
	testhelper.ReadEnvFile(t, "../../")
	logger := zaptest.NewLogger(t)
	ghClient, err := gogithub.NewGQLClient(context.Background(), logger, nil)
	if err != nil {
		t.Skip("skipping test because we can't create a github client")
	}
	wh := NewWorkflow(ghClient, testhelper.EnvOrSkip(t, "WORKFLOW_OWNER"), testhelper.EnvOrSkip(t, "WORKFLOW_REPO"), testhelper.EnvOrSkip(t, "WORKFLOW_ID"), testhelper.EnvOrSkip(t, "WORKFLOW_REF"))
	genericNotificationTest(t, wh)
}
