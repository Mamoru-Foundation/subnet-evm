package mamoru

import (
	"math/big"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ava-labs/subnet-evm/mamoru/stats"

	"github.com/Mamoru-Foundation/mamoru-sniffer-go/mamoru_sniffer"
)

type Feeder interface {
	FeedBlock(*types.Block) mamoru_sniffer.Block
	FeedTransactions(blockNumber *big.Int, blockTime uint64, txs types.Transactions, receipts types.Receipts) []mamoru_sniffer.Transaction
	FeedEvents(types.Receipts) []mamoru_sniffer.Event
	FeedCallTraces([]*CallFrame, uint64) []mamoru_sniffer.CallTrace
	Stats() stats.Stats
}
