package client

import (
	"github.com/go-chassis/go-chassis-apm/common"
)

type ApmClient interface {
	CreateSpans(sc *common.SpanContext) ([]interface{}, error)
	EndSpans(spans []interface{}, status int) error
	CreateEntrySpan(sc *common.SpanContext) (interface{}, error)
	CreateExitSpan(sc *common.SpanContext) (interface{}, error)
	EndSpan(sp interface{}, statusCode int) error
}

