package mamoru

import (
	"encoding/json"
	"math/big"
	"strings"
	"sync/atomic"

	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ethereum/go-ethereum/common"
)

type CallFrame struct {
	Type         string
	From         string
	To           string
	Value        uint64
	Gas          uint64
	GasUsed      uint64
	Input        []byte
	Output       string
	Error        string
	RevertReason string
	Depth        uint32
	TxIndex      uint32
	Logs         []callLog `json:"logs,omitempty" rlp:"optional"`
}

type CallTracer struct {
	env       *vm.EVM
	callstack []CallFrame
	config    CallTracerConfig
	interrupt uint32 // Atomic flag to signal execution interruption
	reason    error  // Textual reason for the interruption
	depth     uint32
	txIndex   uint32
}

type CallTracerConfig struct {
	OnlyTopCall bool `json:"onlyTopCall"` // If true, call tracer won't collect any subcalls
	WithLog     bool `json:"withLog"`     // If true, call tracer will collect event logs
}

var _ vm.EVMLogger = &CallTracer{}

// NewCallTracer returns a native go tracer which tracks
// call frames of a tx, and implements vm.EVMLogger.
func NewCallTracer(OnlyTopCall bool) *CallTracer {
	// First callframe contains tx context info
	// and is populated on start and end.
	return &CallTracer{
		callstack: []CallFrame{},
		config:    CallTracerConfig{OnlyTopCall: OnlyTopCall}}
}

// CaptureStart implements the EVMLogger interface to initialize the tracing operation.
func (t *CallTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	t.env = env
	atomic.StoreUint32(&t.depth, 0)
	t.callstack = append(t.callstack, CallFrame{
		Type:  "CALL",
		From:  addrToHex(from),
		To:    addrToHex(to),
		Input: copyBytes(input),
		Gas:   gas,
		Value: value.Uint64(),
		Depth: atomic.LoadUint32(&t.depth),
	})
	if create {
		t.callstack[0].Type = "CREATE"
	}
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (t *CallTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	atomic.StoreUint32(&t.depth, 0)
	t.callstack[0].GasUsed = gasUsed
	if err != nil {
		t.callstack[0].Error = err.Error()
		if err.Error() == "execution reverted" && len(output) > 0 {
			t.callstack[0].Output = bytesToHex(output)
		}
	} else {
		t.callstack[0].Output = bytesToHex(output)
	}
}

// CaptureState implements the EVMLogger interface to trace a single step of VM execution.
func (t *CallTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

// CaptureFault implements the EVMLogger interface to trace an execution fault.
func (t *CallTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, _ *vm.ScopeContext, depth int, err error) {
}

// CaptureEnter is called when EVM enters a new scope (via call, create or selfdestruct).
func (t *CallTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	if t.config.OnlyTopCall {
		return
	}
	// Skip if tracing was interrupted
	if atomic.LoadUint32(&t.interrupt) > 0 {
		t.env.Cancel()
		return
	}
	atomic.AddUint32(&t.depth, 1)
	var valueC uint64
	if value != nil {
		valueC = value.Uint64()
	}

	call := CallFrame{
		Type:    typ.String(),
		From:    addrToHex(from),
		To:      addrToHex(to),
		Input:   copyBytes(input),
		Gas:     gas,
		Depth:   atomic.LoadUint32(&t.depth),
		Value:   valueC,
		TxIndex: atomic.LoadUint32(&t.txIndex),
	}
	t.callstack = append(t.callstack, call)
}

// CaptureExit is called when EVM exits a scope, even if the scope didn't
// execute any code.
func (t *CallTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	atomic.AddUint32(&t.depth, ^uint32(0))
}

func (*CallTracer) CaptureTxStart(uint64) {}

func (t *CallTracer) CaptureTxEnd(uint64) {
	atomic.AddUint32(&t.txIndex, 1)
}

func (t *CallTracer) GetResult() (json.RawMessage, error) {
	callstack, err := t.TakeResult()
	if err != nil {
		return nil, err
	}

	res, err := MarshalJSON(callstack)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(res), t.reason
}

// TakeResult returns the json-encoded nested list of call traces, and any
// error arising from the encoding or forceful termination (via `Stop`).
func (t *CallTracer) TakeResult() ([]*CallFrame, error) {
	var frames []*CallFrame
	for _, call := range t.callstack {
		rcall := call
		frames = append(frames, &rcall)
	}

	defer func() {
		t.callstack = []CallFrame{}
		atomic.StoreUint32(&t.interrupt, 0)
		atomic.StoreUint32(&t.txIndex, 0)
		t.reason = nil
	}()

	return frames, t.reason
}

// Stop terminates execution of the tracer at the first opportune moment.
func (t *CallTracer) Stop(err error) {
	t.reason = err
	atomic.StoreUint32(&t.interrupt, 1)
}

func bytesToHex(s []byte) string {
	return "0x" + common.Bytes2Hex(s)
}

func addrToHex(a common.Address) string {
	return strings.ToLower(a.Hex())
}

func copyBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)

	return c
}
