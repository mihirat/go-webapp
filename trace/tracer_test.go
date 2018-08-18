package trace

import (
	"bytes"
	"testing"
)

func TestOff(t *testing.T) {
	var silenceTracer Tracer = Off()
	silenceTracer.Trace("data")
}
func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("new returns nil")
	} else {
		tracer.Trace("hello, trace")
		if buf.String() != "hello, trace\n" {
			t.Errorf("wrong message as '%s'", buf.String())
		}
	}
}
