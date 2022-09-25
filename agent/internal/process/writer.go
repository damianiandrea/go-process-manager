package process

import (
	"context"

	"github.com/damianiandrea/go-process-manager/agent/internal/message"
)

type OutputWriter struct {
	ctx         context.Context
	processUuid string

	msgProducer message.ProcessOutputMsgProducer
}

func NewOutputWriter(ctx context.Context, processUuid string, msgProducer message.ProcessOutputMsgProducer) *OutputWriter {
	return &OutputWriter{ctx: ctx, processUuid: processUuid, msgProducer: msgProducer}
}

func (w *OutputWriter) Write(p []byte) (int, error) {
	output := &message.ProcessOutput{ProcessUuid: w.processUuid, Data: p}
	err := w.msgProducer.Produce(w.ctx, output)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
