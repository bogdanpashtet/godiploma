package cipher

import cipherv1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/cipher/v1"

type Method string

const (
	MethodUnspecified Method = "unspecified"
	MethodLSB         Method = "lsb"
)

func (t Method) FromDomain() cipherv1.Method {
	switch t {
	case MethodLSB:
		return cipherv1.Method_METHOD_LSB
	default:
		return cipherv1.Method_METHOD_UNSPECIFIED
	}
}

func ConvertMethodToDomain(t cipherv1.Method) Method {
	switch t {
	case cipherv1.Method_METHOD_LSB:
		return MethodLSB
	default:
		return MethodUnspecified
	}
}
