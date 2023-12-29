package sync_state

import (
	"context"
	"os"
	"testing"

	"github.com/ava-labs/avalanchego/utils/rpc"
	"github.com/stretchr/testify/assert"
)

// MockInfoClient is a mock implementation of the InfoClient interface
type MockInfoClient struct {
	Bootstrapped bool
	Err          error
}

func (m *MockInfoClient) IsBootstrapped(context.Context, string, ...rpc.Option) (bool, error) {
	return m.Bootstrapped, m.Err
}

func TestFetchSyncStatus(t *testing.T) {
	mockClient := &MockInfoClient{Bootstrapped: true, Err: nil}
	client := &Client{
		apiClient:  mockClient,
		Url:        "mockURL",
		PolishTime: 1,
		quit:       make(chan struct{}),
		signals:    make(chan os.Signal, 1),
	}

	// Call the method you want to test
	client.fetchSyncStatus()

	// Assert that the isSync field was updated correctly
	assert.True(t, client.isSync.Load(), "Expected isSync to be true")
}

func TestFetchNotSyncStatus(t *testing.T) {
	mockClient := &MockInfoClient{Bootstrapped: false, Err: nil}
	client := &Client{
		apiClient:  mockClient,
		Url:        "mockURL",
		PolishTime: 1,
		quit:       make(chan struct{}),
		signals:    make(chan os.Signal, 1),
	}

	// Call the method you want to test
	client.fetchSyncStatus()

	// Assert that the isSync field was updated correctly
	assert.False(t, client.isSync.Load(), "Expected isSync to be true")
}
