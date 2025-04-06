package v1

import (
	"context"

	cipherd "github.com/bogdanpashtet/godiploma/internal/domain/cipher"
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
)

type Service interface {
	CreateStegoImage(ctx context.Context, req cipherd.CreateStegoImageRequest) ([]filed.File, error)
	Extract(ctx context.Context, req cipherd.ExtractRequest) ([]string, error)
}
