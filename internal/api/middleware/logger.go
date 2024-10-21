package middleware

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type LoggingOption func(l *loggingInterceptor)

func WithLogger(logger *slog.Logger) LoggingOption {
	return func(l *loggingInterceptor) {
		l.logger = logger
	}
}

func Logging(opts ...LoggingOption) connect.Interceptor {
	l := &loggingInterceptor{
		logger: slog.Default(),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

type loggingInterceptor struct {
	logger *slog.Logger
}

func (l loggingInterceptor) WrapUnary(
	unaryFunc connect.UnaryFunc,
) connect.UnaryFunc {
	return func(
		ctx context.Context, req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		if l.logger.Enabled(ctx, slog.LevelDebug) {
			logMsg(ctx, l.logger, req.Any(), recv)
		}
		res, err := unaryFunc(ctx, req)
		if err == nil && l.logger.Enabled(ctx, slog.LevelDebug) {
			logMsg(ctx, l.logger, res.Any(), send)
		}
		l.log(ctx, req.Spec(), err)
		return res, err
	}
}

var marshaller = protojson.MarshalOptions{UseProtoNames: true}

const (
	recv = "RECV: "
	send = "SEND: "
)

func logMsg(ctx context.Context, logger *slog.Logger, v any, prefix string) {
	b, err := marshaller.Marshal(v.(proto.Message))
	if err != nil {
		panic(err)
	}
	if len(b) > 4096 {
		b = append(b[:4096], "..."...)
	}
	logger.DebugContext(ctx, prefix+string(b))
}

func (l loggingInterceptor) WrapStreamingClient(
	clientFunc connect.StreamingClientFunc,
) connect.StreamingClientFunc {
	return clientFunc
}

func (l loggingInterceptor) WrapStreamingHandler(
	handlerFunc connect.StreamingHandlerFunc,
) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context, conn connect.StreamingHandlerConn,
	) error {
		err := handlerFunc(ctx, conn)
		l.log(ctx, conn.Spec(), err)
		return err
	}
}

func (l loggingInterceptor) log(
	ctx context.Context, spec connect.Spec, err error,
) {
	lvl := slog.LevelInfo
	msg := spec.Procedure + " "

	if err == nil {
		msg += "ok"
	} else {
		var ce *connect.Error
		if !errors.As(err, &ce) {
			ce = connect.NewError(connect.CodeUnknown, err)
		}
		if ce.Code() == connect.CodeUnknown {
			lvl = slog.LevelError
		}
		msg += ce.Error()
	}

	l.logger.Log(ctx, lvl, msg)
}
