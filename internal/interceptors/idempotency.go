package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const ContextHeaderKey = "Idempotency-Key"

type idempotencyKey struct{}

func AddIdempotencyKeyToCtx(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, idempotencyKey{}, key)
}

func ExtractIdempotencyKeyFromCtx(ctx context.Context) string {
	if key, ok := ctx.Value(idempotencyKey{}).(string); ok {
		return key
	}

	return ""
}

func IdempotencyUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		headers := md.Get(ContextHeaderKey)
		if len(headers) > 0 {
			// To log idempotency key in kibana ServerInterceptor should be placed before grpcx.LoggingServerInterceptor.
			ctx = zap.AddToCtx(ctx, zap.String("idempotency-key", headers[0]))
			ctx = AddIdempotencyKeyToCtx(ctx, headers[0])
		}

		return handler(ctx, req)
	}
}
