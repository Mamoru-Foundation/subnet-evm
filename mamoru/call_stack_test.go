package mamoru

import (
	"github.com/ava-labs/subnet-evm/core/vm"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"

	"math/big"
	"testing"
)

func TestCallStackTracer_CaptureEnter(t1 *testing.T) {
	t := NewCallStackTracer(nil, RandStr(3), false, "test")

	//Start Tx
	t.CaptureTxStart(0)
	// depth 0
	t.CaptureStart(nil, common.Address{}, common.Address{}, false, nil, 0, big.NewInt(0))
	// depth 1
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	// depth 2
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	// depth 3
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	// depth 4
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	t.CaptureExit([]byte{}, 0, nil) // depth 4
	t.CaptureExit([]byte{}, 0, nil) // depth 3
	t.CaptureExit([]byte{}, 0, nil) // depth 2
	t.CaptureExit([]byte{}, 0, nil) // depth 1

	t.CaptureEnd(nil, 0, nil) // depth 0
	t.CaptureTxEnd(0)

	t.CaptureTxStart(0)
	t.CaptureStart(nil, common.Address{}, common.Address{}, false, nil, 0, big.NewInt(0))
	t.CaptureEnd(nil, 0, nil)
	t.CaptureTxEnd(0)

	result, err := t.TakeResult()
	assert.NoError(t1, err)
	assert.NotNil(t1, result)
	assert.Equal(t1, 6, len(result))

	assert.Equal(t1, uint32(0), result[0].Depth)
	assert.Equal(t1, uint32(0), result[0].TxIndex)

	assert.Equal(t1, uint32(1), result[1].Depth)

	assert.Equal(t1, uint32(2), result[2].Depth)
	assert.Equal(t1, uint32(0), result[2].TxIndex)

	assert.Equal(t1, uint32(3), result[3].Depth)
	assert.Equal(t1, uint32(0), result[3].TxIndex)

	assert.Equal(t1, uint32(4), result[4].Depth)
	assert.Equal(t1, uint32(0), result[4].TxIndex)

	assert.Equal(t1, uint32(0), result[5].Depth)
	assert.Equal(t1, uint32(1), result[5].TxIndex)
}

func TestToFlatten(t1 *testing.T) {
	frames := []CallStackFrame{
		{Depth: 0,
			Calls: []CallStackFrame{
				{Depth: 1,
					Calls: []CallStackFrame{
						{Depth: 2, Calls: []CallStackFrame{}},
					},
				},
			},
		},
	}
	want := []CallStackFrame{
		{Depth: 0},
		{Depth: 1},
		{Depth: 2},
	}

	have := toFlatten(frames)
	for i := range have {
		assert.Equal(t1, want[i].Depth, have[i].Depth)
	}
}
