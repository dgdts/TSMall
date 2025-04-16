package xjsoniter

import (
	"unsafe"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type ignoreOmitEmptyTagExtension struct {
	jsoniter.DummyExtension
}

type ignoreOmitEmptyTagEncoder struct {
	originDecoder jsoniter.ValEncoder
}

// key logic
func (p *ignoreOmitEmptyTagEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

func (p *ignoreOmitEmptyTagEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	p.originDecoder.Encode(ptr, stream)
}

func (e *ignoreOmitEmptyTagExtension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return &ignoreOmitEmptyTagEncoder{encoder}
}

func init() {
	hlog.Debugf("fix omitemptytag issue")
	jsoniter.RegisterExtension(&ignoreOmitEmptyTagExtension{})
}
