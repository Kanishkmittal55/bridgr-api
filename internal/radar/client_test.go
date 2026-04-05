package radar

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	t.Run("creates client with default addr", func(t *testing.T) {
		// Use an invalid addr - dial will fail but we verify client struct is created
		// In real scenario, localhost:50051 would work when radar is running
		client, err := NewClient(ctx, Config{
			Addr:         "localhost:50051",
			CrawlTimeout: 5 * time.Minute,
		})
		// Connection may succeed if nothing is listening (connection refused)
		// or succeed if radar is running
		if err != nil {
			t.Skipf("radar not available: %v", err)
			return
		}
		require.NotNil(t, client)
		defer client.Close()
	})

	t.Run("uses default timeout when zero", func(t *testing.T) {
		client, err := NewClient(ctx, Config{
			Addr:         "localhost:50051",
			CrawlTimeout: 0,
		})
		if err != nil {
			t.Skipf("radar not available: %v", err)
			return
		}
		require.NotNil(t, client)
		assert.Equal(t, DefaultCrawlTimeout, client.crawlTimeout)
		defer client.Close()
	})

	t.Run("uses empty addr default", func(t *testing.T) {
		client, err := NewClient(ctx, Config{})
		if err != nil {
			t.Skipf("radar not available: %v", err)
			return
		}
		require.NotNil(t, client)
		defer client.Close()
	})
}

func TestConfig_Defaults(t *testing.T) {
	assert.Equal(t, 30*time.Minute, DefaultCrawlTimeout)
}
