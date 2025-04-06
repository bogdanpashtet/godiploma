package v1

import (
	cipherd "github.com/bogdanpashtet/godiploma/internal/domain/cipher"
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	cipherv1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/cipher/v1"
	"github.com/samber/lo"
)

func convertCreateStegoImageToDomain(req *cipherv1.CreateStegoImageRequest) cipherd.CreateStegoImageRequest {
	return cipherd.CreateStegoImageRequest{
		Method:    cipherd.ConvertMethodToDomain(req.Method),
		Plaintext: req.Plaintext,
		Files: lo.Map(req.Files, func(file *cipherv1.File, _ int) filed.File {
			return filed.File{
				Metadata: filed.Metadata{
					Type: filed.ConvertTypeToDomain(file.Metadata.Type),
				},
				File: file.DocumentData,
			}
		}),
	}
}

func convertExtractToDomain(req *cipherv1.ExtractRequest) cipherd.ExtractRequest {
	return cipherd.ExtractRequest{
		Method: cipherd.ConvertMethodToDomain(req.Method),
		Files: lo.Map(req.Files, func(file *cipherv1.File, _ int) filed.File {
			return filed.File{
				Metadata: filed.Metadata{
					Type: filed.ConvertTypeToDomain(file.Metadata.Type),
				},
				File: file.DocumentData,
			}
		}),
	}
}
