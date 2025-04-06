package interceptors

import (
	"context"
	"time"

	"github.com/bogdanpashtet/godiploma/internal/log"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	grpcServiceKey         = "grpc.service"
	grpcMethodKey          = "grpc.method"
	grpcCodeKey            = "grpc.code"
	grpcDurationKey        = "grpc.duration"
	grpcClientNamespaceKey = "grpc.client_namespace"
)

func LoggingUnaryServerInterceptor(logger log.Logger) grpc.UnaryServerInterceptor {
	disableMethods, processedFields := processProtoFieldsAndMethods()

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := info.Server.(healthpb.HealthServer); ok {
			return handler(ctx, req)
		}
		startTime := time.Now()

		fields := getCallFields(ctx, info.FullMethod, startTime)

		// Safe to ignore existence, as we only read from map and metadata
		// should always be present.
		md, _ := metadata.FromIncomingContext(ctx)
		if namespace, ok := clientNamespaceFromMetadata(md); ok {
			fields = append(fields, brlog.String(grpcClientNamespaceKey, namespace))
		}

		reqCtx := brlog.AddToCtx(ctx, fields...)

		methodFullName := toMethodFullName(info.FullMethod)
		isLogAvailable := true
		if _, ok := disableMethods[methodFullName]; ok {
			isLogAvailable = false
		}

		var loggedReq interface{}
		if isLogAvailable {
			loggedReq = req
		}

		writeLog(reqCtx, logger, methodFullName, "server request", grpcRequestKey, loggedReq, processedFields)

		resp, err := handler(reqCtx, req)

		respCtx := brlog.AddToCtx(
			reqCtx,
			brlog.Duration(grpcDurationKey, time.Since(startTime)),
			brlog.Bool(grpcIsLogAvailableKey, isLogAvailable),
			brlog.String(grpcCodeKey, status.Code(err).String()),
			brlog.Error(err),
		)

		var loggedResp interface{}
		if isLogAvailable {
			loggedResp = resp
		}

		writeLog(respCtx, logger, methodFullName, "server response", grpcResponseKey, loggedResp, processedFields)

		return resp, err
	}
}
