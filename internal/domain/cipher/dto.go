package cipher

import filed "github.com/bogdanpashtet/godiploma/internal/domain/file"

type CreateStegoImageRequest struct {
	Method    Method
	Plaintext string
	Files     []filed.File
}
