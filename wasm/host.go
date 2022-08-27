//go:build !tinygo.wasm

// This file is designed to be imported by hosts.

package wasm

import (
	"context"
	"errors"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

func ReadMemory(ctx context.Context, m api.Module, offset, size uint32) ([]byte, error) {
	buf, ok := m.Memory().Read(ctx, offset, size)
	if !ok {
		return nil, fmt.Errorf("Memory.Read(%d, %d) out of range", offset, size)
	}
	return buf, nil
}

func WriteMemory(ctx context.Context, m api.Module, data []byte) (uint64, error) {
	malloc := m.ExportedFunction("malloc")
	if malloc == nil {
		return 0, errors.New("malloc is not exported")
	}
	results, err := malloc.Call(ctx, uint64(len(data)))
	if err != nil {
		return 0, err
	}
	dataPtr := results[0]

	// The pointer is a linear memory offset, which is where we write the name.
	if !m.Memory().Write(ctx, uint32(dataPtr), data) {
		return 0, fmt.Errorf("Memory.Write(%d, %d) out of range of memory size %d",
			dataPtr, len(data), m.Memory().Size(ctx))
	}

	return dataPtr, nil
}