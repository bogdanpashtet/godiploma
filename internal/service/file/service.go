package file

import (
	"context"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"go.uber.org/zap"
)

const kind = "FileService"

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

func (s *Service) UploadFiles(ctx context.Context, req filed.UploadFilesRequest) error {
	s.l = s.l.With(zap.String("kind", kind))

	for i, doc := range req.Files {
		docType := doc
		docData := doc.Type

		docLogger := s.l.With(
			zap.Any("documentIndex", i),
			zap.Any("documentType", docType),
			zap.Any("documentSize", len(docData)),
		)

		docLogger.Info("Processing document")
	}

	return nil
}
