package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/bogdanpashtet/godiploma/internal/config"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.uber.org/zap"
)

const basicAuthScheme = "Basic "

type Authenticator struct {
	cfg    config.AuthConfig
	logger *zap.Logger
}

func NewAuthenticator(cfg *config.AppConfig, logger *zap.Logger) (*Authenticator, error) {
	if len(cfg.Auth.Keys) == 0 {
		logger.Warn("Authenticator initialized without keys. Authentication will be effectively disabled or fail.")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger cannot be nil")
	}
	return &Authenticator{
		cfg:    cfg.Auth,
		logger: logger,
	}, nil
}

func (a *Authenticator) Authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
	}

	authHeaders := md.Get("authorization") // Basic Auth использует заголовок authorization
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header is not provided")
	}
	authHeader := authHeaders[0]

	// Проверяем схему Basic
	if !strings.HasPrefix(authHeader, basicAuthScheme) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization scheme: expected Basic")
	}

	// Декодируем Base64
	encodedCredentials := strings.TrimPrefix(authHeader, basicAuthScheme)
	credentialsBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		a.logger.Warn("Failed to decode base64 basic auth credentials", zap.Error(err))
		return nil, status.Errorf(codes.Unauthenticated, "invalid base64 encoding in authorization header")
	}
	credentials := string(credentialsBytes)

	// Разделяем "username:password"
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		a.logger.Warn("Invalid basic auth credential format (missing colon)", zap.String("credentials", credentials))
		return nil, status.Errorf(codes.Unauthenticated, "invalid credential format")
	}
	username := parts[0]
	password := parts[1]

	// Ищем пользователя и его хэш
	storedHash, found := a.cfg.Keys[username]
	if !found {
		a.logger.Warn("Basic auth user not found", zap.String("username", username))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials") // Общая ошибка
	}

	// Сравниваем хэш из конфига с присланным паролем
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		// Если ошибка - это несовпадение пароля, просто возвращаем Unauthenticated
		if err == bcrypt.ErrMismatchedHashAndPassword {
			a.logger.Warn("Basic auth password mismatch", zap.String("username", username))
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}
		// Другая ошибка bcrypt (маловероятно, но возможно)
		a.logger.Error("Error comparing basic auth password hash", zap.String("username", username), zap.Error(err))
		return nil, status.Error(codes.Internal, "authentication error")
	}

	// Аутентификация успешна!
	a.logger.Debug("Basic authentication successful", zap.String("username", username))

	return ctx, nil
}
