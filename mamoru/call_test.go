package mamoru

import (
	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"math/big"
	"testing"
)

func TestCallTracer_CaptureEnter(t1 *testing.T) {
	t := NewCallTracer(false)

	// Depth 0
	t.CaptureStart(nil, common.Address{}, common.Address{}, false, nil, 0, big.NewInt(0))

	// Depth 1
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)

	// Depth 2
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	t.CaptureExit([]byte{}, 0, nil)

	t.CaptureExit([]byte{}, 0, nil)

	// Depth 1
	t.CaptureEnter(vm.CALL, common.Address{}, common.Address{}, []byte{}, 0, nil)
	t.CaptureExit([]byte{}, 0, nil)

	t.CaptureEnd(nil, 0, nil)

	result, err := t.TakeResult()
	if err != nil {
		return
	}

	assert.Equal(t1, uint32(0), result[0].Depth)
	assert.Equal(t1, uint32(1), result[1].Depth)
	assert.Equal(t1, uint32(2), result[2].Depth)
	assert.Equal(t1, uint32(1), result[3].Depth)
}
