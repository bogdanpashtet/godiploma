package file

import filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"

type Type string

const (
	TypeUnspecified   Type = ""
	TypePassport      Type = "passport"
	TypeDriverLicense Type = "driver_license"
	TypeBankStatement Type = "bank_statement"
)

func (t Type) FromDomain() filev1.Type {
	switch t {
	case TypePassport:
		return filev1.Type_TYPE_PASSPORT
	case TypeDriverLicense:
		return filev1.Type_TYPE_DRIVER_LICENSE
	case TypeBankStatement:
		return filev1.Type_TYPE_BANK_STATEMENT
	default:
		return filev1.Type_TYPE_UNSPECIFIED
	}
}

func ConvertTypeToDomain(t filev1.Type) Type {
	switch t {
	case filev1.Type_TYPE_PASSPORT:
		return TypePassport
	case filev1.Type_TYPE_DRIVER_LICENSE:
		return TypeDriverLicense
	case filev1.Type_TYPE_BANK_STATEMENT:
		return TypeBankStatement
	default:
		return TypeUnspecified
	}
}
