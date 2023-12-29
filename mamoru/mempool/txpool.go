package mempool

import (
	"github.com/ava-labs/subnet-evm/core"
	"github.com/ethereum/go-ethereum/event"
)

type BcTxPool interface {
	// SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription
	//SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
}
