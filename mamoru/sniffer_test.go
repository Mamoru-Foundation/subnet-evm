package mamoru

import (
	"fmt"

	"os"
	"testing"

	"github.com/Mamoru-Foundation/mamoru-sniffer-go/mamoru_sniffer"

	"github.com/stretchr/testify/assert"
)

func TestSniffer_isSnifferEnable(t *testing.T) {
	t.Run("TRUE env is set 1", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "1")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.True(t, got)
	})
	t.Run("TRUE env is set true", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.True(t, got)
	})
	t.Run("FALSE env is set 0", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "0")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
	t.Run("FALSE env is set 0", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "false")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
	t.Run("FALSE env is not set", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "")
		defer unsetEnvSnifferEnable()
		got := isSnifferEnable()
		assert.False(t, got)
	})
}

func unsetEnvSnifferEnable() {
	_ = os.Unsetenv("MAMORU_SNIFFER_ENABLE")
}

func TestSniffer_connect(t *testing.T) {
	t.Run("TRUE ", func(t *testing.T) {
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		got := connect()
		assert.True(t, got)
	})
	t.Run("FALSE connect have error", func(t *testing.T) {
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, fmt.Errorf("Some err") }
		got := connect()
		assert.False(t, got)
	})
}

type MocSyncProcess struct {
	isSync bool
}

func newMocSyncProcess(isSync bool) *MocSyncProcess {
	return &MocSyncProcess{isSync: isSync}
}
func (m *MocSyncProcess) Process() bool {
	return m.isSync
}

func TestSniffer_CheckRequirements1(t *testing.T) {
	t.Run("TRUE ", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			synced: true,
		}
		assert.True(t, s.CheckRequirements())
	})
	t.Run("FALSE chain not sync ", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			synced: true,
			status: newMocSyncProcess(false),
		}
		assert.False(t, s.CheckRequirements())
	})
	t.Run("FALSE connect error", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "true")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, fmt.Errorf("Some err") }
		s := &Sniffer{
			synced: true,
		}
		assert.False(t, s.CheckRequirements())
	})
	t.Run("FALSE env not set", func(t *testing.T) {
		_ = os.Setenv("MAMORU_SNIFFER_ENABLE", "0")
		defer unsetEnvSnifferEnable()
		SnifferConnectFunc = func() (*mamoru_sniffer.Sniffer, error) { return nil, nil }
		s := &Sniffer{
			synced: false,
			status: newMocSyncProcess(false),
		}
		assert.False(t, s.CheckRequirements())
	})
}

func Test_getDeltaBlocks(t *testing.T) {
	t.Run("Success return delta from env ", func(t *testing.T) {
		want := int64(100)
		t.Setenv("MAMORU_SNIFFER_DELTA", fmt.Sprintf("%d", want))

		got := getDelta()
		assert.Equal(t, want, got)
	})
	t.Run("Success return delta from env ", func(t *testing.T) {
		want := DefaultDelta
		got := getDelta()
		assert.Equal(t, want, got)
	})
}
