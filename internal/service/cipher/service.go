package cipher

import (
	"context"
	"fmt"

	cipherd "github.com/bogdanpashtet/godiploma/internal/domain/cipher"
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"go.uber.org/zap"
)

const kind = "CipherService"

type Service struct {
	l *zap.Logger
}

func New(
	l *zap.Logger,
) *Service {
	return &Service{
		l: l,
	}
}

func (s *Service) CreateStegoImage(ctx context.Context, req cipherd.CreateStegoImageRequest) ([]filed.File, error) {
	s.l = s.l.With(zap.String("kind", kind))

	resFiles := make([]filed.File, 0, len(req.Files))
	for _, file := range req.Files {
		switch req.Method {
		case cipherd.MethodLSB:
			res, err := cipherLSB(ctx, req.Plaintext, file)
			if err != nil {
				return nil, fmt.Errorf("fail to cipherLSB: %w", err)
			}
			resFiles = append(resFiles, res)
		default:
		}
	}

	return resFiles, nil
}

func (s *Service) Extract(ctx context.Context, req cipherd.ExtractRequest) ([]string, error) {
	s.l = s.l.With(zap.String("kind", kind))

	resPlaintexts := make([]string, 0, len(req.Files))
	for _, file := range req.Files {
		switch req.Method {
		case cipherd.MethodLSB:
			res, err := extractLSB(ctx, file.File)
			if err != nil {
				return nil, fmt.Errorf("fail to cipherLSB: %w", err)
			}
			resPlaintexts = append(resPlaintexts, res)
		default:
		}
	}

	return resPlaintexts, nil
}
