package file

import cipherv1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/cipher/v1"

type Type string

const (
	TypeUnspecified Type = ""
	TypeBMP         Type = "bmp"
	TypePNG         Type = "png"
)

func (t Type) FromDomain() cipherv1.Type {
	switch t {
	case TypeBMP:
		return cipherv1.Type_TYPE_BMP
	case TypePNG:
		return cipherv1.Type_TYPE_PNG
	default:
		return cipherv1.Type_TYPE_UNSPECIFIED
	}
}

func ConvertTypeToDomain(t cipherv1.Type) Type {
	switch t {
	case cipherv1.Type_TYPE_BMP:
		return TypeBMP
	case cipherv1.Type_TYPE_PNG:
		return TypePNG
	default:
		return TypeUnspecified
	}
}
