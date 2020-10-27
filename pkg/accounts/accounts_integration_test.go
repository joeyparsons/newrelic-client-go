// +build integration

package accounts

import (
	"testing"

	"github.com/stretchr/testify/require"

	mock "github.com/joeyparsons/newrelic-client-go/pkg/testhelpers"
)

func TestIntegrationAccounts(t *testing.T) {
	t.Parallel()

	client := newIntegrationTestClient(t)

	// Test: List
	params := ListAccountsParams{
		Scope: &RegionScopeTypes.GLOBAL,
	}

	accounts, err := client.ListAccounts(params)
	require.NoError(t, err)
	require.Greater(t, len(accounts), 0)
}

func newIntegrationTestClient(t *testing.T) Accounts {
	tc := mock.NewIntegrationTestConfig(t)

	return New(tc)
}
