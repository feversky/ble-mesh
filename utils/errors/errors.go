package errors

import (
	"fmt"
)

const (
	InvalidKey = ErrorType(iota)
	NilPointer
	InvalidNetKeyIndex
	InvalidAppKeyIndex
	NetKeyNotFoundByNid
	FailedToGenerateNonce
	NoValidNetKeyForDecryption
	NoValidAppKeyForDecryption
	TransportSarFailed
	NodeAddressNotFound
	Timeout
	DataLengthCheckFailed
	NotAPointer
	NetKeyNotBindedToNode
	NetKeyAlreadyBindedToNode
	OpcodeNotSupportedByNode
	AppKeyNotBindedToNode
	AppKeyAlreadyBindedToNode
	NoNetKeyBindedToAppKey
	ElementNotFound
	ModelNotFound
	NotFound
	ModelAppKeyBindingNotFound
	AppKeyAlreadyBindedToModel
	NoAppKeyBindedToModel
	AddressNotInSubscriptionList
	AddressAlreadyInSubscriptionList
	WrongTTLSetting
	WrongGattProxySetting
	InvalidResponse

	CannotSetRangeMin
	CannotSetRangeMax

	//BitString
	WrongFormatOfBitString
	LengthMismatchOfBitString
	OutLenMismatchOfBitString
	BitLenOver64OfBitString
	FieldOfStructNotAddressable
)

var errorDict = map[ErrorType]string{
	InvalidKey:                       "",
	NilPointer:                       "nil pointer",
	InvalidNetKeyIndex:               "net key index is invalid",
	InvalidAppKeyIndex:               "app key index is invalid",
	NetKeyNotFoundByNid:              "net key not found by nid",
	FailedToGenerateNonce:            "failed to generate nonce",
	NoValidNetKeyForDecryption:       "no valid net key was found for decryption",
	NoValidAppKeyForDecryption:       "no valid app key was found for decryption",
	TransportSarFailed:               "failed to reassembly transport message",
	NodeAddressNotFound:              "node address not found",
	NotAPointer:                      "not a pointer to struct",
	NetKeyNotBindedToNode:            "netkey is not binded to the node yet",
	NetKeyAlreadyBindedToNode:        "netkey is already binded to the node",
	OpcodeNotSupportedByNode:         "the model does not contain this opcode",
	AppKeyNotBindedToNode:            "appkey is not binded to the node yet",
	AppKeyAlreadyBindedToNode:        "appkey is already binded to the node",
	NoNetKeyBindedToAppKey:           "no netkey is binded with this appkey",
	ElementNotFound:                  "element not found",
	ModelNotFound:                    "model not found",
	ModelAppKeyBindingNotFound:       "cannot find a model binded with this appkey",
	AppKeyAlreadyBindedToModel:       "appkey is already binded to this model",
	NoAppKeyBindedToModel:            "no appkey is binded to this model yet",
	AddressNotInSubscriptionList:     "address not in subscription list",
	AddressAlreadyInSubscriptionList: "address already in subscription list",
	WrongTTLSetting:                  "wrong ttl setting",
	WrongGattProxySetting:            "wrong gatt proxy setting",
	InvalidResponse:                  "invalid response, request is failed",
	Timeout:                          "timeout happens",
	DataLengthCheckFailed:            "data length check failed",

	CannotSetRangeMin: "Cannot Set Range Min",
	CannotSetRangeMax: "Cannot Set Range Max",

	WrongFormatOfBitString:      "format of BitString is wrong",
	LengthMismatchOfBitString:   "length of data does not match the bitstring when unpacking",
	OutLenMismatchOfBitString:   "length of outputs does not match the bitstring",
	BitLenOver64OfBitString:     "bit length greater than 64 is not supported",
	FieldOfStructNotAddressable: "filed of struct is not addressable",
}

type ErrorType uint

type MeshError struct {
	errorType ErrorType
	context   []interface{}
}

func (et ErrorType) New() *MeshError {
	return &MeshError{errorType: et, context: []interface{}{}}
}

func (err *MeshError) AddContext(c interface{}) error {
	err.context = append(err.context, c)
	return err
}

func (err *MeshError) AddContextF(format string, args ...interface{}) error {
	err.context = append(err.context, fmt.Sprintf(format, args...))
	return err
}

func (err *MeshError) Error() string {
	if len(err.context) > 0 {
		return fmt.Sprintf("error type: %s, context: %+#v)", errorDict[err.errorType], err.context)
	}
	return errorDict[err.errorType]
}

func (err *MeshError) ErrorType() ErrorType {
	return err.errorType
}
