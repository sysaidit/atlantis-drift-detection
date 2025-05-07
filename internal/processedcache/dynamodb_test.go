package processedcache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sysaidit/atlantis-drift-detection/internal/testhelper"
)

func makeTestClient(t *testing.T) *DynamoDB {
	testhelper.ReadEnvFile(t, "../../")
	client, err := NewDynamoDB(context.Background(), testhelper.EnvOrSkip(t, "DYNAMODB_TABLE"))
	require.NoError(t, err)
	return client
}

func TestDynamoDB(t *testing.T) {
	GenericCacheWorkflowTest(t, makeTestClient(t))
}
