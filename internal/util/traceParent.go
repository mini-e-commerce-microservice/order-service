package util

import (
	"context"
	whttp "github.com/SyaibanAhmadRamadhan/http-wrapper"
)

func GetTraceParent(ctx context.Context) *string {
	traceParent := whttp.GetTraceParent(ctx)
	if traceParent != "" {
		return &traceParent
	}

	return nil
}
