package mamoru

import (
	"math/big"
	"math/rand"
	"sync"
	"time"

	"github.com/ava-labs/subnet-evm/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Mamoru-Foundation/mamoru-sniffer-go/mamoru_sniffer"
)

const (
	CtxBlockchain = "blockchain"
	CtxTxpool     = "txpool"
)

type Tracer struct {
	feeder  Feeder
	mu      sync.Mutex
	builder mamoru_sniffer.EvmCtxBuilder
}

func NewTracer(feeder Feeder) *Tracer {
	builder := mamoru_sniffer.NewEvmCtxBuilder()
	tr := &Tracer{builder: builder, feeder: feeder}
	return tr
}

// SetTxpoolCtx Set the context of the sniffer
func (t *Tracer) SetTxpoolCtx() {
	t.builder.SetMempoolSource()
}

func (t *Tracer) SetStatistics(blocks, txs, events, callTraces uint64) {
	t.builder.SetStatistics(blocks, txs, events, callTraces)
}

func (t *Tracer) FeedBlock(block *types.Block) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.builder.SetBlock(
		t.feeder.FeedBlock(block),
	)
}

func (t *Tracer) FeedTransactions(blockNumber *big.Int, blockTime uint64, txs types.Transactions, receipts types.Receipts) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.builder.AppendTxs(
		t.feeder.FeedTransactions(blockNumber, blockTime, txs, receipts),
	)
}

func (t *Tracer) FeedEvents(receipts types.Receipts) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.builder.AppendEvents(
		t.feeder.FeedEvents(receipts),
	)
}

func (t *Tracer) FeedCallTraces(callFrames []*CallFrame, blockNumber uint64) {
	defer t.mu.Unlock()
	t.mu.Lock()
	t.builder.AppendCallTraces(
		t.feeder.FeedCallTraces(callFrames, blockNumber),
	)
}

func (t *Tracer) Send(start time.Time, blockNumber *big.Int, blockHash common.Hash, snifferContext string) {
	defer t.mu.Unlock()
	t.mu.Lock()

	if sniffer != nil {
		stats := t.feeder.Stats()
		defer stats.Reset()
		log.Info("Mamoru Stats",
			"blocks", stats.GetBlocks(),
			"txs", stats.GetTxs(),
			"events", stats.GetEvents(),
			"callTraces", stats.GetTraces(),
			"lastBlockNum", blockNumber.String(),
			"lastBlockHash", blockHash.TerminalString(),
			"ctx", snifferContext)

		t.builder.SetStatistics(stats.GetBlocks(), stats.GetTxs(), stats.GetEvents(), stats.GetTraces())
		t.builder.SetBlockData(blockNumber.String(), blockHash.String())
		sniffer.ObserveEvmData(t.builder.Finish())
	}

	logCtx := []interface{}{
		"elapsed", common.PrettyDuration(time.Since(start)),
		"number", blockNumber.String(),
		"hash", blockHash,
		"ctx", snifferContext,
	}
	log.Info("Mamoru Sniffer finish", logCtx...)
}

func RandStr(length int) string {
	rand.NewSource(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
