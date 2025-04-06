package v1

import (
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
	"github.com/samber/lo"
)

func convertUploadRequestToDomain(req *filev1.UploadFilesRequest) filed.UploadFilesRequest {
	return filed.UploadFilesRequest{Files: lo.Map(req.Documents, func(file *filev1.File, _ int) filed.File {
		return filed.File{
			Type: filed.ConvertTypeToDomain(file.DocumentType),
			File: file.DocumentData,
		}
	})}
}
