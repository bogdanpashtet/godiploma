package v1

import (
	"context"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
)

type Service interface {
	UploadFiles(ctx context.Context, req filed.UploadFilesRequest) error
}
