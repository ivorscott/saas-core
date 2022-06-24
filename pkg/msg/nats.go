package msg

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"syscall"

	"github.com/devpies/saas-core/pkg/web"

	"go.uber.org/zap"
)

// StreamContext manages the StreamContext.
type StreamContext struct {
	logger   *zap.Logger
	shutdown chan os.Signal
	js       nats.JetStreamContext
}

// NewStreamContext returns a new StreamContext.
func NewStreamContext(logger *zap.Logger, shutdown chan os.Signal, address string, port string) *StreamContext {
	var err error

	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", address, port))
	if err != nil {
		logger.Error("could not connect to NATS", zap.Error(err))
	}

	js, err := nc.JetStream()
	if err != nil {
		logger.Error("could not create JetStream context", zap.Error(err))
	}

	return &StreamContext{logger, shutdown, js}
}

// Create creates the named stream.
func (jctx *StreamContext) Create(streamName string) *nats.StreamInfo {
	strInfo, err := jctx.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{fmt.Sprintf("%s.*", streamName), streamName},
		MaxAge:   0, // Keep forever.
		Storage:  nats.FileStorage,
	})
	if err != nil {
		jctx.logger.Error("could not create stream", zap.Error(err))
	}
	return strInfo
}

// Publish publishes a message.
func (jctx *StreamContext) Publish(subject string, message []byte) {
	ack, err := jctx.js.Publish(subject, message)
	if err != nil {
		jctx.logger.Info("failed to publish message", zap.Error(err))
	}
	jctx.logger.Info(fmt.Sprintf("%v", ack))
}

type listenHandlerFunc func(ctx context.Context, message interface{}) error

func (jctx *StreamContext) Listen(messageType string, subject, queueGroup string, handler listenHandlerFunc, opts ...nats.SubOpt) *nats.Subscription {
	fn := jctx.setupMsgHandler(messageType, handler)
	s, err := jctx.js.QueueSubscribe(subject, queueGroup, fn, opts...)
	if err != nil {
		jctx.logger.Info("subscription failed", zap.Error(err), zap.Any("data", jctx.js))
	}
	return s
}

func (jctx *StreamContext) setupMsgHandler(messageType string, handler listenHandlerFunc) func(msg *nats.Msg) {
	return func(m *nats.Msg) {
		message, err := UnmarshalMsg(m.Data)
		if err != nil {
			jctx.logger.Error("error decoding message", zap.Error(err), zap.String("message", string(m.Data)))
			return
		}

		if string(message.Type) != messageType {
			jctx.logger.Info(
				"warning message type was not expected",
				zap.String("wanted", messageType),
				zap.String("got", string(message.Type)),
			)
			return
		}

		jctx.logger.Info(
			"processing message",
			zap.String("type", messageType),
			zap.String("tenantID", message.Metadata.TenantID),
			zap.String("traceID", message.Metadata.TraceID),
		)
		err = handler(addMetadataContext(message.Metadata), &message)

		switch err.(type) {
		case nil:
			err = m.Ack()
			if err != nil {
				jctx.logger.Error("error acknowledging message", zap.Error(err))
			}
		case *web.Error:
			jctx.logger.Error("error handling message", zap.Error(err))
		case *web.Shutdown:
			jctx.logger.Error("integrity issue: shutting down service", zap.Error(err))
			jctx.shutdown <- syscall.SIGSTOP
		default:
			panic(err)
		}
	}
}

func addMetadataContext(metadata Metadata) context.Context {
	ctx := context.Background()
	values := web.Values{
		TenantID: metadata.TenantID,
		UserID:   metadata.UserID,
		TraceID:  metadata.TraceID,
	}
	return web.NewContext(ctx, &values)
}
