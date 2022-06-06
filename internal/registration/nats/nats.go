package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type JetStreamContext struct {
	logger *zap.Logger
	js     nats.JetStreamContext
}

// NewJetStreamContext returns a new JetStreamContext.
func NewJetStreamContext(logger *zap.Logger, address string, port string) *JetStreamContext {
	var err error
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", address, port))
	if err != nil {
		logger.Error("could not connect to NATS", zap.Error(err))
	}

	js, err := nc.JetStream()
	if err != nil {
		logger.Error("could not create JetStream context", zap.Error(err))
	}

	return &JetStreamContext{logger, js}
}

// Create creates the named stream.
func (jctx *JetStreamContext) Create(streamName string) *nats.StreamInfo {
	strInfo, err := jctx.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{fmt.Sprintf("%s.*", streamName), streamName},
		MaxAge:   0, // 0 means keep forever
		Storage:  nats.FileStorage,
	})
	if err != nil {
		jctx.logger.Error("could not create stream", zap.Error(err))
	}
	return strInfo
}

// Publish publishes a message.
func (jctx *JetStreamContext) Publish(subject string, message []byte) {
	ack, err := jctx.js.Publish(subject, message)
	if err != nil {
		jctx.logger.Info("failed to publish message", zap.Error(err))
	}
	jctx.logger.Info(fmt.Sprintf("%v", ack))
}
