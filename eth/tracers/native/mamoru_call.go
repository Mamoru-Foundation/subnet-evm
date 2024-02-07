package native

import (
	"encoding/json"

	"github.com/ava-labs/subnet-evm/core/vm"
	"github.com/ava-labs/subnet-evm/eth/tracers"
	"github.com/ava-labs/subnet-evm/mamoru"
)

func init() {
	tracers.DefaultDirectory.Register("mamoru", newMamoruCallStackTracer, false)
}

var _ vm.EVMLogger = (*mamoru.CallTracer)(nil)
var _ tracers.Tracer = (*mamoru.CallTracer)(nil)

var _ vm.EVMLogger = (*mamoru.CallStackTracer)(nil)
var _ tracers.Tracer = (*mamoru.CallStackTracer)(nil)

// Used to register the tracer with the tracer manager
func newMamoruCallStackTracer(ctx *tracers.Context, cfg json.RawMessage) (tracers.Tracer, error) {
	var config mamoru.CallTracerConfig
	if cfg != nil {
		if err := json.Unmarshal(cfg, &config); err != nil {
			return nil, err
		}
	}
	return &mamoru.CallStackTracer{
		Source: "console",
		Config: config,
	}, nil
}
