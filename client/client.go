package client

import (
	"github.com/go-chassis/go-chassis-apm/common"
)

//ApmClient for apm interface
type ApmClient interface {
	CreateEntrySpan(sc *common.SpanContext) (interface{}, error)
	CreateExitSpan(sc *common.SpanContext) (interface{}, error)
	EndSpan(sp interface{}, statusCode int) error
}
