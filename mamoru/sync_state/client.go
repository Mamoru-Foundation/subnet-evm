package sync_state

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/utils/rpc"

	"github.com/ethereum/go-ethereum/log"
)

type InfoClient interface {
	IsBootstrapped(context.Context, string, ...rpc.Option) (bool, error)
}

type Client struct {
	apiClient  InfoClient
	isSync     atomic.Bool
	Url        string
	PolishTime uint

	quit    chan struct{}
	signals chan os.Signal
}

func NewJSONRPCRequest(url string, PolishTime uint) *Client {
	c := &Client{
		apiClient:  info.NewClient(url),
		Url:        url,
		PolishTime: PolishTime,
		quit:       make(chan struct{}),
		signals:    make(chan os.Signal, 1),
	}
	c.isSync.Store(false)
	// Register for SIGINT (Ctrl+C) and SIGTERM (kill) signals
	signal.Notify(c.signals, syscall.SIGINT, syscall.SIGTERM)

	return c
}

func (c *Client) loop() {
	ticker := time.NewTicker(time.Duration(c.PolishTime) * time.Second)
	defer ticker.Stop()
	// Perform the first tick immediately
	c.fetchSyncStatus()

	for {
		select {
		case <-ticker.C:
			c.fetchSyncStatus()
		case <-c.quit:
			log.Info("Mamoru SyncProcess Shutting down...")
			return
		case <-c.signals:
			log.Info("Signal received, initiating shutdown...")
			c.Close()
		}
	}
}

func (c *Client) fetchSyncStatus() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("Mamoru requesting syncData ...")
	// Send the request and get the response
	isBootstrapped, err := c.apiClient.IsBootstrapped(ctx, "X")
	if err != nil {
		log.Error("Mamoru Sync", "err", err)
		return
	}

	c.isSync.Store(isBootstrapped)
}

func (c *Client) Close() {
	close(c.quit)
}
