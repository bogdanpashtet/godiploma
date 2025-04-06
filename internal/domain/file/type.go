package file

import filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"

type Type string

const (
	TypeUnspecified   Type = ""
	TypePassport      Type = "passport"
	TypeDriverLicense Type = "driver_license"
	TypeBankStatement Type = "bank_statement"
)

func (t Type) FromGrpc() filev1.DocumentType {
	switch t {
	case TypePassport:
		return filev1.DocumentType_DOCUMENT_TYPE_PASSPORT
	case TypeDriverLicense:
		return filev1.DocumentType_DOCUMENT_TYPE_DRIVER_LICENSE
	case TypeBankStatement:
		return filev1.DocumentType_DOCUMENT_TYPE_BANK_STATEMENT
	default:
		return filev1.DocumentType_DOCUMENT_TYPE_UNSPECIFIED
	}
}
