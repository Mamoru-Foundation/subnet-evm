package sync_state

import (
	"github.com/ethereum/go-ethereum/log"
)

type SyncProcess struct {
	client *Client
}

func NewSyncProcess(opNodeRpcUrl string, PolishTime uint) *SyncProcess {
	return &SyncProcess{
		client: NewJSONRPCRequest(opNodeRpcUrl, PolishTime),
	}
}

func (s SyncProcess) Start() {
	go s.client.loop()
}

func (s SyncProcess) Stop() {
	s.client.Close()
}

func (s SyncProcess) Process() bool {
	isSync := s.client.isSync.Load()
	log.Info("Mamoru SyncProcess", "sync", isSync)

	return isSync
}
